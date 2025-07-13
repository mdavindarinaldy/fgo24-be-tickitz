package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

func AddPaymentMethod(newData dto.NewData) error {
	if newData.Name == "" {
		return errors.New("payment method name should not be empty")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(
		context.Background(),
		`
		INSERT INTO payment_methods 
			(name, created_at)
		VALUES
			($1,$2)`, newData.Name, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func GetPaymentMethod() ([]dto.SubData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.SubData{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT id, name FROM payment_methods`)
	if err != nil {
		return []dto.SubData{}, err
	}
	genres, err := pgx.CollectRows[dto.SubData](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.SubData{}, err
	}
	return genres, nil
}

func AddTransactions(newTrx dto.NewTrx, userId int) (int, int, error) {
	if newTrx.Cinema == "" || newTrx.Location == "" || newTrx.Date == "" || newTrx.Showtime == "" || newTrx.PaymentMethodId == 0 || newTrx.MovieId == 0 || len(newTrx.Seats) == 0 {
		return 0, 0, errors.New("transactions data should not be empty")
	}

	showtimeParsed, err := time.Parse("15:04:05", newTrx.Showtime)
	if err != nil {
		return 0, 0, errors.New("invalid showtime format, use HH:MM:SS")
	}

	dateParsed, err := time.Parse("2006-01-02", newTrx.Date)
	if err != nil {
		return 0, 0, errors.New("invalid date format, use YYYY-MM-DD")
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	var showtimeId int
	err = tx.QueryRow(context.Background(),
		`SELECT id FROM showtimes 
		WHERE 
			id_movie = $1 
			AND cinema = $2 
			AND location = $3 
			AND date = $4 
			AND showtime = $5`,
		newTrx.MovieId, newTrx.Cinema, newTrx.Location, dateParsed, showtimeParsed).Scan(&showtimeId)

	if err == pgx.ErrNoRows {
		err = tx.QueryRow(context.Background(),
			`INSERT INTO showtimes 
			 (id_movie, cinema, location, date, showtime)
			 VALUES ($1, $2, $3, $4, $5) 
			 RETURNING id`,
			newTrx.MovieId, newTrx.Cinema, newTrx.Location, dateParsed, showtimeParsed).Scan(&showtimeId)
		if err != nil {
			return 0, 0, err
		}
	} else if err != nil {
		return 0, 0, err
	}

	var price float64
	err = tx.QueryRow(context.Background(),
		`SELECT price FROM movies WHERE id = $1`,
		newTrx.MovieId).Scan(&price)
	if err != nil {
		return 0, 0, err
	}

	totalAmount := price * float64(len(newTrx.Seats))

	var transactionId int
	err = tx.QueryRow(context.Background(), `
		INSERT INTO transactions 
		(id_user, id_payment_method, total_amount)
		VALUES ($1, $2, $3) RETURNING id`,
		userId, newTrx.PaymentMethodId, totalAmount).Scan(&transactionId)
	if err != nil {
		return 0, 0, err
	}

	for _, seat := range newTrx.Seats {
		_, err = tx.Exec(context.Background(),
			`INSERT INTO transactions_detail 
			(id_transaction, id_showtime, seat)
			 VALUES ($1, $2, $3)`,
			transactionId, showtimeId, seat)
		if err != nil {
			return 0, 0, err
		}
	}

	return showtimeId, transactionId, nil
}

func GetTransactionsHistory(userId int) ([]dto.TransactionHistory, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.TransactionHistory{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT 
			m.id AS id_movie, m.title, s.location, 
    		s.cinema, s.showtime, s.date, s.id AS id_showtime, 
			t.id AS id_transaction, string_agg(td.seat,', ') AS seats 
			FROM transactions t
		JOIN transactions_detail td ON td.id_transaction = t.id
		JOIN showtimes s ON s.id = td.id_showtime
		JOIN movies m ON m.id = s.id_movie
		WHERE t.id_user=$1
		GROUP BY m.id, t.id, s.id
		ORDER BY t.created_at DESC;`, userId)
	if err != nil {
		return []dto.TransactionHistory{}, err
	}

	// trxHistory, err := pgx.CollectRows[dto.TransactionHistory](rows, pgx.RowToStructByName)
	// if err != nil {
	// 	return []dto.TransactionHistory{}, err
	// }

	var trxHistory []dto.TransactionHistory
	for rows.Next() {
		var currentTrx dto.TrxHistory
		err := rows.Scan(&currentTrx.MovieId, &currentTrx.MovieTitle, &currentTrx.Location, &currentTrx.Cinema, &currentTrx.Showtime, &currentTrx.Date, &currentTrx.ShowtimeId, &currentTrx.TransactionId, &currentTrx.Seats)
		if err != nil {
			return nil, err
		}
		trxHistory = append(trxHistory, dto.TransactionHistory{
			MovieId:       currentTrx.MovieId,
			MovieTitle:    currentTrx.MovieTitle,
			Location:      currentTrx.Location,
			Cinema:        currentTrx.Cinema,
			Showtime:      dto.CustomTime{Time: currentTrx.Showtime},
			Date:          dto.CustomDate{Time: currentTrx.Date},
			ShowtimeId:    currentTrx.ShowtimeId,
			TransactionId: currentTrx.TransactionId,
			Seats:         currentTrx.Seats,
		})
	}

	return trxHistory, nil
}

func GetReservedSeat(req dto.ReservedSeatsRequest) (dto.ReservedSeatsResponse, error) {
	if req.MovieId == 0 || req.Cinema == "" || req.Location == "" || req.Date == "" || req.Showtime == "" {
		return dto.ReservedSeatsResponse{}, errors.New("all fields must be provided")
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return dto.ReservedSeatsResponse{}, err
	}
	defer conn.Close()

	var showtimeID int
	err = conn.QueryRow(context.Background(),
		`SELECT id FROM showtimes 
         WHERE id_movie = $1 
         AND cinema = $2 
         AND location = $3 
         AND date = $4 
         AND showtime = $5`,
		req.MovieId, req.Cinema, req.Location, req.Date, req.Showtime).Scan(&showtimeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ReservedSeatsResponse{
				ShowtimeId: 0,
				Seats:      "",
			}, nil
		}
		return dto.ReservedSeatsResponse{}, err
	}

	row, err := conn.Query(context.Background(), `
		SELECT id_showtime, string_agg(seat, ', ') AS seats 
		FROM transactions_detail 
		WHERE id_showtime = $1
		GROUP BY id_showtime`, showtimeID)
	if err != nil {
		return dto.ReservedSeatsResponse{}, err
	}

	var reservedSeats dto.ReservedSeatsResponse
	reservedSeats, err = pgx.CollectOneRow[dto.ReservedSeatsResponse](row, pgx.RowToStructByName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ReservedSeatsResponse{
				ShowtimeId: showtimeID,
				Seats:      "",
			}, nil
		}
		return dto.ReservedSeatsResponse{}, err
	}

	return reservedSeats, nil
}

func GetSalesPerMovie() ([]dto.SalesPerMovie, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []dto.SalesPerMovie{}, err
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(),
		`
		SELECT m.id AS id_movie, m.title, 
		COUNT(td.seat) AS tickets_sold, 
		m.price AS price_per_ticket,
		COUNT(td.seat)*m.price AS total_amount 
		FROM transactions_detail td
		JOIN showtimes s ON s.id=td.id_showtime
		JOIN movies m ON m.id=s.id_movie
		GROUP BY m.id;
		`)

	if err != nil {
		return []dto.SalesPerMovie{}, err
	}

	data, err := pgx.CollectRows[dto.SalesPerMovie](rows, pgx.RowToStructByName)
	if err != nil {
		return []dto.SalesPerMovie{}, err
	}

	return data, nil
}
