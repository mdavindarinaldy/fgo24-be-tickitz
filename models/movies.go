package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetMovies(search string, filter string, page int) ([]dto.Movie, utils.PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return nil, utils.PageData{}, err
	}
	defer conn.Close()

	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	var totalData int
	err = conn.QueryRow(context.Background(), `
		SELECT COUNT(DISTINCT m.id)
		FROM movies m
		JOIN movies_genres mg ON mg.id_movie = m.id
		JOIN genres g ON g.id = mg.id_genre
		WHERE m.title ILIKE $1 AND g.name ILIKE $2
	`, "%"+search+"%", "%"+filter+"%").Scan(&totalData)
	if err != nil {
		return nil, utils.PageData{}, err
	}

	totalPage := (totalData + limit - 1) / limit

	rows, err := conn.Query(context.Background(), `
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
			m.id, m.title, m.synopsis, m.release_date, m.price, m.runtime,
			m.poster, m.backdrop,
			g.genres, d.directors, c.casts
		FROM movies m
		LEFT JOIN genres_agg g ON m.id = g.id_movie
		LEFT JOIN directors_agg d ON m.id = d.id_movie
		LEFT JOIN casts_agg c ON m.id = c.id_movie
		WHERE m.title ILIKE $1 AND g.genres ILIKE $2
		OFFSET $3 LIMIT $4
	`, "%"+search+"%", "%"+filter+"%", offset, limit)
	if err != nil {
		return nil, utils.PageData{}, err
	}

	movies, err := pgx.CollectRows[dto.Movie](rows, pgx.RowToStructByName)
	if err != nil {
		return nil, utils.PageData{}, err
	}

	pageData := utils.PageData{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: page,
	}

	return movies, pageData, nil
}

func GetMovie(id int) (dto.Movie, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return dto.Movie{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(),
		`WITH genres_agg AS (
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
			m.id, 
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
		WHERE id = $1`, id)
	if err != nil {
		return dto.Movie{}, err
	}
	movie, err := pgx.CollectOneRow[dto.Movie](rows, pgx.RowToStructByName)
	if err != nil {
		return dto.Movie{}, err
	}
	return movie, nil
}

func GetUpcomingMovies() ([]dto.Movie, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.Movie{}, err
	}
	defer conn.Close()
	time := time.Now()
	rows, err := conn.Query(context.Background(),
		`WITH genres_agg AS (
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
			m.id, 
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
		WHERE release_date>$1`, time)
	if err != nil {
		return []dto.Movie{}, err
	}
	movies, err := pgx.CollectRows[dto.Movie](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.Movie{}, err
	}
	return movies, nil
}

func GetGenre() ([]dto.SubData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.SubData{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT id, name FROM genres`)
	if err != nil {
		return []dto.SubData{}, err
	}
	genres, err := pgx.CollectRows[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.SubData{}, err
	}
	return genres, nil
}

func GetDirectors(search string) ([]dto.SubData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.SubData{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT id, name FROM directors WHERE name ILIKE $1`, "%"+search+"%")
	if err != nil {
		return []dto.SubData{}, err
	}
	genres, err := pgx.CollectRows[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.SubData{}, err
	}
	return genres, nil
}

func GetCasts(search string) ([]dto.SubData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.SubData{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT id, name FROM casts WHERE name ILIKE $1`, "%"+search+"%")
	if err != nil {
		return []dto.SubData{}, err
	}
	genres, err := pgx.CollectRows[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.SubData{}, err
	}
	return genres, nil
}

func AddMovie(newMovie dto.NewMovie, adminId int) error {
	if newMovie.Title == "" || newMovie.Synopsis == "" || newMovie.ReleaseDate == "" || newMovie.Price == 0 || newMovie.Runtime == 0 || newMovie.Poster == nil || newMovie.Backdrop == nil || newMovie.Genres == "" || newMovie.Directors == "" || newMovie.Casts == "" {
		return errors.New("new movie data should not be empty")
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	var newMovieId int

	err = tx.QueryRow(
		context.Background(),
		`
		INSERT INTO movies 
			(created_by, title, synopsis, release_date, price, runtime, poster, backdrop, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;
		`, adminId, newMovie.Title, newMovie.Synopsis, newMovie.ReleaseDate, newMovie.Price,
		newMovie.Runtime, newMovie.Poster, newMovie.Backdrop, time.Now()).Scan(&newMovieId)
	if err != nil {
		return err
	}

	directors := strings.Split(newMovie.Directors, ", ")
	for _, v := range directors {
		directorId, _ := strconv.Atoi(v)
		_, err = tx.Exec(context.Background(),
			`
			INSERT INTO movies_directors
				(id_director, id_movie)
			VALUES
				($1, $2);
			`, directorId, newMovieId)
		if err != nil {
			return err
		}
	}

	casts := strings.Split(newMovie.Casts, ", ")
	for _, v := range casts {
		castId, _ := strconv.Atoi(v)
		_, err = tx.Exec(context.Background(),
			`
			INSERT INTO movies_casts
				(id_cast, id_movie)
			VALUES
				($1, $2);
			`, castId, newMovieId)
		if err != nil {
			return err
		}
	}

	genres := strings.Split(newMovie.Genres, ", ")
	for _, v := range genres {
		genreId, _ := strconv.Atoi(v)
		_, err = tx.Exec(context.Background(),
			`
			INSERT INTO movies_genres
				(id_genre, id_movie)
			VALUES
				($1, $2);
			`, genreId, newMovieId)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddDirector(data string) (dto.SubData, error) {
	if data == "" {
		return dto.SubData{}, errors.New("director name should not be empty")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return dto.SubData{}, err
	}
	defer conn.Close()

	rows, err := conn.Query(
		context.Background(),
		`
		INSERT INTO directors 
			(name, created_at)
		VALUES
			($1,$2)
		RETURNING id, name;
		`, data, time.Now())
	if err != nil {
		return dto.SubData{}, err
	}

	result, err := pgx.CollectOneRow[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return dto.SubData{}, err
	}

	return result, nil
}

func AddCast(data string) (dto.SubData, error) {
	if data == "" {
		return dto.SubData{}, errors.New("cast name should not be empty")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return dto.SubData{}, err
	}
	defer conn.Close()

	rows, err := conn.Query(
		context.Background(),
		`
		INSERT INTO casts 
			(name, created_at)
		VALUES
			($1, $2)
		RETURNING id, name;
		`, data, time.Now())
	if err != nil {
		return dto.SubData{}, err
	}

	result, err := pgx.CollectOneRow[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return dto.SubData{}, err
	}

	return result, nil
}

func AddGenre(data string) (dto.SubData, error) {
	if data == "" {
		return dto.SubData{}, errors.New("genre name should not be empty")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return dto.SubData{}, err
	}
	defer conn.Close()

	rows, err := conn.Query(
		context.Background(),
		`
		INSERT INTO genres 
			(name, created_at)
		VALUES
			($1,$2)
		RETURNING id, name;
		`, data, time.Now())
	if err != nil {
		return dto.SubData{}, err
	}

	result, err := pgx.CollectOneRow[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return dto.SubData{}, err
	}

	return result, nil
}

func UpdateMovie(updateMovie dto.NewMovie, movieId int) error {
	if updateMovie.Title == "" || updateMovie.Synopsis == "" || updateMovie.ReleaseDate == "" || updateMovie.Price == 0 || updateMovie.Runtime == 0 || updateMovie.Poster == nil || updateMovie.Backdrop == nil || updateMovie.Genres == "" || updateMovie.Directors == "" || updateMovie.Casts == "" {
		return errors.New("update movie data should not be empty")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	_, err = tx.Exec(
		context.Background(),
		`UPDATE movies SET title = $1, synopsis = $2, release_date = $3,
		price = $4, runtime = $5, poster = $6, backdrop = $7
		WHERE id = $8`,
		updateMovie.Title, updateMovie.Synopsis, updateMovie.ReleaseDate,
		updateMovie.Price, updateMovie.Runtime, updateMovie.Poster, updateMovie.Backdrop, movieId)

	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(),
		`DELETE FROM movies_genres WHERE id_movie = $1`,
		movieId)
	if err != nil {
		return err
	}

	genres := strings.Split(updateMovie.Genres, ", ")
	for _, newVal := range genres {
		genreId, _ := strconv.Atoi(newVal)
		_, err = tx.Exec(context.Background(),
			`
			INSERT INTO movies_genres
				(id_genre, id_movie)
			VALUES
				($1, $2);
			`, genreId, movieId)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(context.Background(),
		`DELETE FROM movies_directors WHERE id_movie = $1`,
		movieId)
	if err != nil {
		return err
	}

	directors := strings.Split(updateMovie.Directors, ", ")
	for _, newVal := range directors {
		directorId, _ := strconv.Atoi(newVal)
		_, err = tx.Exec(context.Background(),
			`
			INSERT INTO movies_directors
				(id_director, id_movie)
			VALUES
				($1, $2);
			`, directorId, movieId)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(context.Background(),
		`DELETE FROM movies_casts WHERE id_movie = $1`,
		movieId)
	if err != nil {
		return err
	}

	casts := strings.Split(updateMovie.Casts, ", ")
	for _, newVal := range casts {
		castId, _ := strconv.Atoi(newVal)
		_, err = tx.Exec(context.Background(),
			`
			INSERT INTO movies_casts
				(id_cast, id_movie)
			VALUES
				($1, $2);
			`, castId, movieId)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteMovie(movieId int) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(context.Background(), `DELETE FROM movies WHERE id = $1`, movieId)
	if err != nil {
		return err
	}
	return nil
}
