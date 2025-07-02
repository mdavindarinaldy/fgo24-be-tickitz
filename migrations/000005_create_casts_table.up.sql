CREATE TABLE casts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE movies_casts (
    id_cast INT REFERENCES casts(id),
    id_movie INT REFERENCES movies(id),
    created_at TIMESTAMP DEFAULT NOW()
);