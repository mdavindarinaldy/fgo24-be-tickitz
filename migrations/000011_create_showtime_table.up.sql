CREATE TABLE showtimes (
    id SERIAL PRIMARY KEY,
    id_movie INT REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE,
    cinema VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    showtime TIME NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_showtime UNIQUE (id_movie, cinema, location, date, showtime)
);

ALTER TABLE transactions_detail
ADD COLUMN id_showtime INT;

ALTER TABLE transactions_detail
ADD CONSTRAINT transactions_detail_id_showtime_fkey
    FOREIGN KEY (id_showtime) REFERENCES showtimes(id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT unique_seat_per_showtime UNIQUE (id_showtime, seat);

ALTER TABLE transactions
DROP COLUMN cinema,
DROP COLUMN location,
DROP COLUMN showtime,
DROP CONSTRAINT transactions_id_movie_fkey,
DROP COLUMN id_movie;