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