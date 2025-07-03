ALTER TABLE movies_genres
DROP CONSTRAINT movies_genres_pkey;

ALTER TABLE movies_genres
DROP CONSTRAINT movies_genres_id_genre_fkey,
DROP CONSTRAINT movies_genres_id_movie_fkey;

ALTER TABLE movies_genres
ADD CONSTRAINT movies_genres_id_genre_fkey
    FOREIGN KEY (id_genre) REFERENCES genres(id),
ADD CONSTRAINT movies_genres_id_movie_fkey
    FOREIGN KEY (id_movie) REFERENCES movies(id);