package dto

import (
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.IsZero() {
		return []byte(`"0000-00-00"`), nil
	}
	return []byte(`"` + cd.Format("2006-01-02") + `"`), nil
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		cd.Time = time.Time{}
		return nil
	}
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.IsZero() {
		return []byte(`"0000-00-00"`), nil
	}
	return []byte(`"` + ct.Format("15:04:05") + `"`), nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		ct.Time = time.Time{}
		return nil
	}
	t, err := time.Parse("15:04:05", str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

type NewTrx struct {
	MovieId         int      `json:"movieId" db:"id_movie"`
	PaymentMethodId int      `json:"paymentMethodId" db:"id_payment_method"`
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
	MovieId       int        `json:"movieId" db:"id_movie" redis:"movieId"`
	MovieTitle    string     `json:"movieTitle" db:"title" redis:"movieTitle"`
	Location      string     `json:"location" db:"location" redis:"location"`
	Cinema        string     `json:"cinema" db:"cinema" redis:"cinema"`
	Date          CustomDate `json:"date" db:"date" redis:"date"`
	Showtime      CustomTime `json:"showtime" db:"showtime" redis:"showtime"`
	Seats         string     `json:"seats" db:"seats" redis:"seats"`
	ShowtimeId    int        `json:"showtimeId" db:"id_showtime" redis:"showtimeId"`
	TransactionId int        `json:"transactionId" db:"id_transaction" redis:"transactionId"`
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

type ReservedSeatsRequest struct {
	MovieId  int    `form:"id_movie" binding:"required"`
	Cinema   string `form:"cinema" binding:"required"`
	Location string `form:"location" binding:"required"`
	Date     string `form:"date" binding:"required"`
	Showtime string `form:"showtime" binding:"required"`
}

type ReservedSeatsResponse struct {
	ShowtimeId int    `json:"id_showtime" db:"id_showtime"`
	Seats      string `json:"seats" db:"seats"`
}

type SalesPerMovie struct {
	MovieId        int     `json:"id_movie" db:"id_movie"`
	MovieTitle     string  `json:"title" db:"title"`
	TicketsSold    int     `json:"ticketsSold" db:"tickets_sold"`
	PricePerTicket float64 `json:"pricePerTicket" db:"price_per_ticket"`
	TotalAmount    float64 `json:"totalAmount" db:"total_amount"`
}
