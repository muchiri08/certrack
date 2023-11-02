CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(200) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    subscribed bool NOT NULL
);