CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    created_by INT REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    synopsis TEXT NOT NULL,
    release_date DATE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    runtime INT NOT NULL,
    poster VARCHAR(255) NOT NULL,
    backdrop VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);