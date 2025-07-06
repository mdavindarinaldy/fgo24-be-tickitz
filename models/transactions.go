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
	if newTrx.Cinema == "" || newTrx.Location == "" || newTrx.Date == "" || newTrx.Showtime == "" || newTrx.PaymentMethodId == 0 || newTrx.TotalAmount == 0 || newTrx.MovieId == 0 || len(newTrx.Seats) == 0 {
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

	var transactionId int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO transactions 
		(id_user, id_payment_method, total_amount)
		 VALUES ($1, $2, $3) RETURNING id`,
		userId, newTrx.PaymentMethodId, newTrx.TotalAmount).Scan(&transactionId)
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
