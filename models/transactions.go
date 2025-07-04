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
