CREATE TABLE IF NOT EXISTS secrets (
    secret_id uuid primary key,
    owner_id  uuid REFERENCES users (user_id) on delete cascade,
    name           varchar(256) not null,
    kind           integer not null,
    metadata       bytea,
    data           bytea not null
);

CREATE UNIQUE INDEX IX_Owner_Name ON secrets (owner_id, name);
