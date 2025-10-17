CREATE TABLE IF NOT EXISTS users (
    user_id uuid primary key,
    email varchar(128) not null unique,
    password_hash text not null
);
