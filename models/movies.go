package models

import (
	"be-tickitz/utils"
	"context"

	"github.com/jackc/pgx/v5"
)

func GetMovies(search string, filter string, page int) ([]utils.ResponseMovie, utils.PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []utils.ResponseMovie{}, utils.PageData{}, err
	}

	type Count struct {
		Count int `db:"count"`
	}
	count, err := conn.Query(context.Background(),
		`SELECT COUNT(*) as count FROM movies m
		JOIN movies_genres mg ON mg.
		WHERE title ILIKE $1`, "%"+search+"%")
	if err != nil {
		return []utils.ResponseMovie{}, utils.PageData{}, err
	}
	countData, err := pgx.CollectOneRow[Count](count, pgx.RowToStructByName)
	if err != nil {
		return []utils.ResponseMovie{}, utils.PageData{}, err
	}

	offset := (page - 1) * 5
	if page == 0 {
		page = 1
	} else if ((page * 5) - countData.Count) < 5 {
		page = 1
	}

	totalPage := 0
	if countData.Count%5 != 0 {
		totalPage = (countData.Count / 5) + 1
	} else {
		totalPage = countData.Count / 5
	}

	pageData := utils.PageData{
		TotalData:   countData.Count,
		TotalPage:   totalPage,
		CurrentPage: page,
	}
}
