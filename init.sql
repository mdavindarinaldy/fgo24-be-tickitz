INSERT INTO users (name, email, phone_number, password, role, created_at)
VALUES
('admin','admin@gmail.com','0812345','admin123','admin','2025-07-02');

SELECT m.title FROM movies m
JOIN movies_genres mg ON mg.id_movie = m.id
JOIN genres g ON g.id = mg.id_genre
WHERE m.title ILIKE '%a%' AND g.name ILIKE '%Adventure%';

WITH genres_agg AS (
    SELECT mg.id_movie, string_agg(DISTINCT g.name, ', ') AS genres
    FROM movies_genres mg
    JOIN genres g ON g.id = mg.id_genre
    GROUP BY mg.id_movie
),
directors_agg AS (
    SELECT md.id_movie, string_agg(DISTINCT d.name, ', ') AS directors
    FROM movies_directors md
    JOIN directors d ON d.id = md.id_director
    GROUP BY md.id_movie
),
casts_agg AS (
    SELECT mc.id_movie, string_agg(DISTINCT c.name, ', ') AS casts
    FROM movies_casts mc
    JOIN casts c ON c.id = mc.id_cast
    GROUP BY mc.id_movie
)
SELECT 
    m.title, 
    m.synopsis, 
    m.release_date, 
    m.price, 
    m.runtime, 
    m.poster, 
    m.backdrop, 
    g.genres, 
    d.directors, 
    c.casts
FROM movies m
LEFT JOIN genres_agg g ON m.id = g.id_movie
LEFT JOIN directors_agg d ON m.id = d.id_movie
LEFT JOIN casts_agg c ON m.id = c.id_movie
WHERE m.title ILIKE '%a%' AND g.genres ILIKE '%Adventure%';

SELECT id, name FROM casts WHERE name ILIKE '%%';

DELETE FROM directors WHERE id=24;

DELETE FROM movies_casts WHERE id_movie>0;

INSERT INTO movies_directors (id_director, id_movie)
VALUES
(1,1),
(2,2),
(3,3),
(4,3),
(5,4),
(6,5),
(7,6),
(8,7),
(9,8),
(10,9),
(11,10),
(12,11),
(13,12),
(14,13),
(15,13),
(16,14),
(17,15),
(18,16),
(19,17),
(20,18),
(21,19),
(22,20),
(23,20);

INSERT INTO movies_genres (id_genre, id_movie)
VALUES
(8,1),
(4,1),
(8,2),
(9,2),
(3,3),
(12,3),
(1,4),
(7,4),
(11,5),
(17,5),
(2,6),
(15,6),
(1,7),
(17,7),
(1,8),
(15,8),
(7,9),
(17,9),
(1,10),
(15,10),
(1,11),
(2,11),
(11,12),
(11,13),
(11,14),
(11,15),
(4,16),
(5,16),
(1,17),
(17,17),
(2,18),
(4,18),
(1,19),
(17,19),
(18,20),
(1,20);

INSERT INTO payment_methods (name, created_at)
VALUES
('OVO',now()),
('GOPAY',now()),
('DANA',now()),
('VISA',now()),
('G-Pay',now()),
('Paypal',now()),
('BCA',now()),
('BRI',now());

SELECT 
    m.id AS id_movie, m.title, s.location, 
    s.cinema, s.date, s.id AS id_showtime, 
    t.id AS id_transactions, t.id_user, string_agg(td.seat,', ') AS seats 
    FROM transactions t
JOIN transactions_detail td ON td.id_transaction = t.id
JOIN showtimes s ON s.id = td.id_showtime
JOIN movies m ON m.id = s.id_movie
WHERE t.id_user=2
GROUP BY m.id, t.id, s.id;

SELECT id_showtime, string_agg(seat,', ') AS seats 
FROM transactions_detail 
WHERE id_showtime = 1
GROUP BY id_showtime;

SELECT * FROM showtimes 
WHERE id_movie = 1 
AND cinema = 'hiflix'
AND location = 'jakarta'
AND date = '2025-07-13' 
AND showtime = '18:30:00';

SELECT m.id AS id_movie, m.title, 
COUNT(td.seat) AS tickets_sold, 
m.price AS price_per_ticket,
COUNT(td.seat)*m.price AS total_amount
FROM transactions_detail td
JOIN showtimes s ON s.id=td.id_showtime
JOIN movies m ON m.id=s.id_movie
GROUP BY m.id;