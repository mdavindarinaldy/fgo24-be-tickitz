package dto

import "time"

type Movie struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Synopsis    string    `json:"synopsis" db:"synopsis"`
	ReleaseDate time.Time `json:"releaseDate" db:"release_date"`
	Price       float64   `json:"price" db:"price"`
	Runtime     int       `json:"runtime" db:"runtime"`
	Poster      string    `json:"poster" db:"poster"`
	Backdrop    string    `json:"backdrop" db:"backdrop"`
	Genres      string    `json:"genres" db:"genres"`
	Directors   string    `json:"directors" db:"directors"`
	Casts       string    `json:"casts" db:"casts"`
}

type NewMovie struct {
	Title       string  `form:"title" json:"title" db:"title"`
	Synopsis    string  `form:"synopsis" json:"synopsis" db:"synopsis"`
	ReleaseDate string  `form:"releaseDate" json:"releaseDate" db:"release_date"`
	Price       float64 `form:"price" json:"price" db:"price"`
	Runtime     int     `form:"runtime" json:"runtime" db:"runtime"`
	Poster      *string `json:"poster"`
	Backdrop    *string `json:"backdrop"`
	Genres      string  `form:"genres" json:"genres" db:"genres"`
	Directors   string  `form:"directors" json:"directors" db:"directors"`
	Casts       string  `form:"casts" json:"casts" db:"casts"`
}

type NewData struct {
	Name string `json:"name"`
}

type SubData struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
