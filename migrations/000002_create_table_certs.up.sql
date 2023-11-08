CREATE TABLE IF NOT EXISTS certs(
    id serial PRIMARY KEY,
    domain VARCHAR(30) NOT NULL,
    issuer VARCHAR(50) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    days_left INTEGER NOT NULL,
    status VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id INTEGER REFERENCES users(id)
);