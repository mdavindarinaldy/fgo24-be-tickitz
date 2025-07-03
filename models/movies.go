package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetMovies(search string, filter string, page int) ([]dto.Movie, utils.PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}
	defer conn.Close()

	type Count struct {
		Title string `db:"title"`
	}

	count, err := conn.Query(context.Background(),
		`
		SELECT m.title FROM movies m
		JOIN movies_genres mg ON mg.id_movie = m.id
		JOIN genres g ON g.id = mg.id_genre
		WHERE m.title ILIKE $1 AND g.name ILIKE $2
		`, "%"+search+"%", "%"+filter+"%")

	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}

	data, err := pgx.CollectRows[Count](count, pgx.RowToStructByName)
	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}

	limit := 10
	offset := (page - 1) * limit
	if page == 0 {
		page = 1
	} else if ((page * limit) - len(data)) < limit {
		page = 1
	}

	totalPage := 0
	if len(data)%limit != 0 {
		totalPage = (len(data) / limit) + 1
	} else {
		totalPage = len(data) / limit
	}

	pageData := utils.PageData{
		TotalData:   len(data),
		TotalPage:   totalPage,
		CurrentPage: page,
	}

	rows, err := conn.Query(context.Background(),
		`
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
		WHERE m.title ILIKE $1 AND g.genres ILIKE $2
		OFFSET $3
		LIMIT $4;
		`, "%"+search+"%", "%"+filter+"%", offset, limit)

	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}
	users, err := pgx.CollectRows[dto.Movie](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}

	return users, pageData, nil
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

func AddMovie(newMovie dto.Movie, adminId int) error {
	if newMovie.Title == "" || newMovie.Synopsis == "" || newMovie.ReleaseDate.IsZero() || newMovie.Price == 0 || newMovie.Runtime == 0 || newMovie.Poster == "" || newMovie.Backdrop == "" || newMovie.Genres == "" || newMovie.Directors == "" || newMovie.Casts == "" {
		return errors.New("new movie data should not be empty")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	// genres := strings.Split(newMovie.Genres, ", ")
	// directors := strings.Split(newMovie.Directors, ", ")
	// casts := strings.Split(newMovie.Casts, ", ")

	_, err = conn.Exec(
		context.Background(),
		`
		INSERT INTO movie 
			(created_by, title, synopsis, 
			release_date, price, runtime, 
			poster, backdrop, created_at)
		VALUES
			($1,$2,$3,$4,$5,$6,"-","-",$7);
		`, adminId, newMovie.Title, newMovie.Synopsis,
		newMovie.ReleaseDate, newMovie.Price,
		newMovie.Runtime, time.Now())
	if err != nil {
		return err
	}

	// _, err = conn.Exec(
	// 	context.Background(),
	// 	`
	// 	INSERT INTO movie
	// 		(created_by, title, synopsis,
	// 		release_date, price, runtime,
	// 		poster, backdrop, created_at)
	// 	VALUES
	// 		($1,$2,$3,$4,$5,$6,"-","-",$7);
	// 	`, adminId, newMovie.Title, newMovie.Synopsis,
	// 	newMovie.ReleaseDate, newMovie.Price,
	// 	newMovie.Runtime, time.Now())
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
	var row dto.SubData
	err = conn.QueryRow(
		context.Background(),
		`
		INSERT INTO directors 
			(name, created_at)
		VALUES
			($1,$2)
		RETURNING id, name;
		`, data, time.Now()).Scan(&row)
	if err != nil {
		return dto.SubData{}, err
	}
	return row, nil
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
	var row dto.SubData
	err = conn.QueryRow(
		context.Background(),
		`
		INSERT INTO casts 
			(name, created_at)
		VALUES
			($1,$2)
		RETURNING id, name;
		`, data, time.Now()).Scan(&row)
	if err != nil {
		return dto.SubData{}, err
	}
	return row, nil
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
	var row dto.SubData
	err = conn.QueryRow(
		context.Background(),
		`
		INSERT INTO genres 
			(name, created_at)
		VALUES
			($1,$2)
		RETURNING id, name;
		`, data, time.Now()).Scan(&row)
	if err != nil {
		return dto.SubData{}, err
	}
	return row, nil
}
