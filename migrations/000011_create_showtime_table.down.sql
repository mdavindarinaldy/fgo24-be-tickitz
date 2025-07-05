ALTER TABLE transactions
ADD COLUMN cinema VARCHAR(255),
ADD COLUMN location VARCHAR(255),
ADD COLUMN showtime TIMESTAMP;

ALTER TABLE transactions_detail
DROP CONSTRAINT unique_seat_per_showtime,
DROP CONSTRAINT transactions_detail_id_showtime_fkey,
DROP COLUMN id_showtime;

DROP TABLE showtimes;