# GophKeeper

GophKeeper — учебный проект менеджера секретов. Сервер хранит секреты (учётные данные, банковские карты, бинарные файлы и т. п.) для аутентифицированных пользователей. Клиент — CLI для регистрации, входа и управления секретами.

Содержимое репозитория включает:
- HTTP-сервер с JWT-аутентификацией
- Хранение данных в PostgreSQL и миграции
- CLI-клиент (cobra + resty)
- Юнит-тесты для сервисного слоя

## Архитектура

Слои:
- internal/app — инициализация приложения: логгер, конфиг, БД, миграции, сервисы
- internal/server — HTTP-слой
  - httsserver — HTTP-сервер и роутер (chi)
  - middleware/security — извлечение userID из JWT
  - handler — обработчики HTTP
  - service — бизнес-логика (AuthService, SecretService)
  - model — DTO и модели домена
- internal/storage — доступ к БД (pgx)
- internal/util — вспомогательные утилиты (JWT)
- pkg/hasher — HMAC-SHA256 хеширование паролей
- pkg/token — простое хранение токена на клиенте (файл)

Данные:
- users: user_id (uuid), email (unique), password_hash
- secrets: secret_id (uuid), owner_id (FK users), name (unique в рамках owner), kind (int), metadata (bytea), data (bytea)

## Конфигурация

internal/config/config.go:
- Параметры читаются из флагов и ENV
- SERVER_ADDRESS (флаг -a) — адрес сервера (по умолчанию localhost:8080)
- DB (флаг -d) — строка подключения к PostgreSQL (по умолчанию postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable)

JWT:
- Секрет и срок жизни токена задаются в internal/app/app.go (SecretKey, tokenExp)

## Подготовка окружения

1) PostgreSQL: создайте базу и выдайте доступ согласно строке подключения.
2) Миграции применяются автоматически при старте сервера (Database.Migrate()).

Пример для локальной разработки:
- Поднимите PostgreSQL локально или через docker
- Проверьте доступность строки подключения из конфига

## Сборка и запуск

Сервер:
- cd GophKeeper
- go build ./cmd/server
- ./server -a localhost:8080 -d "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

Клиент:
- cd GophKeeper
- go build ./cmd/client
- ./client gophkeeper-cli

CLI команды (из ./client):
- Регистрация: ./client auth register -e user@example.com -p secret
- Вход:        ./client auth login -e user@example.com -p secret
- Секреты:
  - ./client secret cred -n name -l login -p pass -m metadata.json
  - ./client secret card -n name -u 4111111111111111 -e 12/30 -v 123 -o HOLDER -m metadata.json
  - ./client secret bin  -n name -f ./path/to/file -m metadata.bin

Токен сохраняется в файл token.txt (pkg/token). Для защищённых запросов CLI включает заголовок Authorization.

## API

Базовый URL: http://localhost:8080

Аутентификация:
- POST /api/signup — регистрация. Тело: {"email":"...","password":"..."}. В ответе заголовок Authorization: Bearer <token> (201/409).
- POST /api/signin — вход. Тело: как выше. В ответе заголовок Authorization: Bearer <token> (200/401/500).

Секреты (требует JWT, middleware/security):
- POST /api/createsecret — создать. Тело: {name, data, metadata, type}. data и metadata — base64 строки. Коды: 201/409/400/401/500.
- PUT  /api/updatesecret — обновить. Тело: {name, data, metadata, type}. Коды: 200/204/400/401/500.
- GET  /api/{name} — получить. Коды: 200/204/400/401/500. Ответ: {name, data, metadata, type}, где data/metadata — base64.

Примечания по данным:
- Внешний API принимает/возвращает base64 в data/metadata.
- В хранилище используются «сырые» байты.

## Разработка

Тесты:
- cd GophKeeper
- go test ./...

Стиль логирования: zap (internal/logger)
Фреймворк HTTP: chi
DB-драйвер: pgx

## Безопасность
- Хеш пароля: HMAC-SHA256 (pkg/hasher)
- JWT: HS256, срок действия задаётся в app.go
- Middleware ��роверяет токен и прокидывает userID в контекст обработчика

## Ограничения и TODO
- Конфигурация TLS/автосертификатов (autocert) настроена шаблонно; для продакшена требуется настройка доменов и среды.
- В обработчиках встречаются общие сообщения об ошибках — улучшить информативность.
- Секреты валидируются минимально — требуется валидация полей и ограничений размеров.
- Улучшить структуру клиентского CLI (сообщения, коды ошибок) и UX.
