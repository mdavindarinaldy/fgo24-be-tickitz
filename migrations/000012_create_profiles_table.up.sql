CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL UNIQUE,
    profile_picture VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO profiles (id_user, name, phone_number, profile_picture, created_at, updated_at)
SELECT id, name, phone_number, NULL, created_at, updated_at
FROM users;

ALTER TABLE users
DROP COLUMN name,
DROP COLUMN phone_number;