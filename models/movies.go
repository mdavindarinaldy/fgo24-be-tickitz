package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"

	"github.com/jackc/pgx/v5"
)

func GetMovies(search string, filter string, page int) ([]dto.Movie, utils.PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}

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

	offset := (page - 1) * 5
	if page == 0 {
		page = 1
	} else if ((page * 5) - len(data)) < 5 {
		page = 1
	}

	totalPage := 0
	if len(data)%5 != 0 {
		totalPage = (len(data) / 5) + 1
	} else {
		totalPage = len(data) / 5
	}

	pageData := utils.PageData{
		TotalData:   len(data),
		TotalPage:   totalPage,
		CurrentPage: page,
	}

	rows, err := conn.Query(context.Background(),
		`SELECT id, name, email, phone_number FROM users 
		WHERE name ILIKE $1 
		OR phone_number ILIKE $1
		OFFSET $2
		LIMIT 5`, "%"+search+"%", offset)

	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}
	users, err := pgx.CollectRows[dto.Movie](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.Movie{}, utils.PageData{}, err
	}

	return users, pageData, nil
}
