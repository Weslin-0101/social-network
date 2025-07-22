CREATE TABLE users (
    id SERIAL,
    username VARCHAR(100) NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX idx_users_email ON users (email);