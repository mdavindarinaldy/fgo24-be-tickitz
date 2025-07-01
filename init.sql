CREATE TYPE role AS ENUM ('user','admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role role NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

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

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE directors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE casts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE movies_genres (
    id_genre INT REFERENCES genres(id),
    id_movies INT REFERENCES movies(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE movies_directors (
    id_director INT REFERENCES directors(id),
    id_movies INT REFERENCES movies(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE movies_casts (
    id_cast INT REFERENCES casts(id),
    id_movies INT REFERENCES movies(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    id_user INT REFERENCES users(id),
    id_movie INT REFERENCES movies(id),
    id_payment_method INT REFERENCES payment_methods(id),
    total_amount DECIMAL(10,2) NOT NULL,
    location VARCHAR(255) NOT NULL,
    cinema VARCHAR(255) NOT NULL,
    showtime DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions_detail (
    id SERIAL PRIMARY KEY,
    id_transaction INT REFERENCES transactions(id),
    seat VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
