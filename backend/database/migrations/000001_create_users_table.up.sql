CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    email VARCHAR(255),
    hash_password VARCHAR(255),
    created_at TIMESTAMP
)