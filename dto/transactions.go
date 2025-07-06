package dto

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
