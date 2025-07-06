package dto

import "time"

type CustomDate struct {
	time.Time
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + cd.Format("2006-01-02") + `"`), nil
}

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + ct.Format("15:04:05") + `"`), nil
}

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
	MovieId       int        `json:"movieId" db:"id_movie"`
	MovieTitle    string     `json:"movieTitle" db:"title"`
	Location      string     `json:"location" db:"location"`
	Cinema        string     `json:"cinema" db:"cinema"`
	Date          CustomDate `json:"date" db:"date"`
	Showtime      CustomTime `json:"showtime" db:"showtime"`
	Seats         string     `json:"seats" db:"seats"`
	ShowtimeId    int        `json:"showtimeId" db:"id_showtime"`
	TransactionId int        `json:"transactionId" db:"id_transaction"`
}

type TrxHistory struct {
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
