ALTER TABLE movies_casts
DROP CONSTRAINT movies_casts_id_cast_fkey,
DROP CONSTRAINT movies_casts_id_movie_fkey;

ALTER TABLE movies_casts
ADD CONSTRAINT movies_casts_id_cast_fkey
    FOREIGN KEY (id_cast) REFERENCES casts(id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT movies_casts_id_movie_fkey
    FOREIGN KEY (id_movie) REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE movies_casts
ADD CONSTRAINT movies_casts_pkey 
PRIMARY KEY (id_cast, id_movie);