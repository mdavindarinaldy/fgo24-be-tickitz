INSERT INTO users (name, email, phone_number, password, role, created_at)
VALUES
('admin','admin@gmail.com','0812345','admin123','admin','2025-07-02');

SELECT m.title, m.synopsis, m.release_date, m.price, m.runtime, m.poster, m.backdrop FROM movies m
JOIN movies_genres mg ON mg.id_movie = m.id
JOIN genres g ON g.id = mg.id_genre
WHERE m.title ILIKE '%a%' AND g.name ILIKE '%Adventure%';

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