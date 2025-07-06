package dto

import "time"

type NewTrx struct {
	MovieId         int      `json:"movieId" db:"id_movie"`
	PaymentMethodId int      `json:"paymentMethodId" db:"id_payment_method"`
	TotalAmount     float64  `json:"totalAmount" db:"total_amount"`
	Location        string   `json:"location" db:"location"`
	Cinema          string   `json:"cinema" db:"cinema"`
	Date            string   `json:"date" db:"date"`
	Showtime        string   `json:"showtime" db:"id_movie"`
	Seats           []string `json:"seats"`
}

type TrxSuccess struct {
	ShowtimeId    int `json:"showtimeId"`
	TransactionId int `json:"transactionId"`
}

type TransactionHistory struct {
	MovieId       int       `json:"movieId" db:"id_movie"`
	MovieTitle    string    `json:"movieTitle" db:"title"`
	Location      string    `json:"location" db:"location"`
	Cinema        string    `json:"cinema" db:"cinema"`
	Date          time.Time `json:"date" db:"date"`
	Showtime      time.Time `json:"showtime" db:"showtime"`
	Seats         string    `json:"seats" db:"seats"`
	ShowtimeId    int       `json:"showtimeId" db:"id_showtime"`
	TransactionId int       `json:"transactionId" db:"id_transaction"`
}
