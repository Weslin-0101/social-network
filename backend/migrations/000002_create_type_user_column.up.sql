ALTER TABLE users
ADD COLUMN type VARCHAR(50) NOT NULL DEFAULT 'user';

ALTER TABLE users
ADD CONSTRAINT users_type_check CHECK (type IN ('user', 'admin'));