package dto

import "time"

type Movie struct {
	Title       string    `json:"title" db:"title"`
	Synopsis    string    `json:"synopsis" db:"synopsis"`
	ReleaseDate time.Time `json:"releaseDate" db:"release_date"`
	Price       float64   `json:"price" db:"price"`
	Runtime     int       `json:"runtime" db:"runtime"`
	Poster      string    `json:"poster" db:"poster"`
	Backdrop    string    `json:"backdrop" db:"backdrop"`
	Genres      string    `json:"genres" db:"genres"`
	Directors   string    `json:"director" db:"director"`
	Casts       string    `json:"casts" db:"casts"`
}
