ALTER TABLE movies_directors
DROP CONSTRAINT movies_directors_pkey;

ALTER TABLE movies_directors
DROP CONSTRAINT movies_directors_id_director_fkey,
DROP CONSTRAINT movies_directors_id_movie_fkey;

ALTER TABLE movies_directors
ADD CONSTRAINT movies_directors_id_director_fkey
    FOREIGN KEY (id_director) REFERENCES directors(id),
ADD CONSTRAINT movies_directors_id_movie_fkey
    FOREIGN KEY (id_movie) REFERENCES movies(id);