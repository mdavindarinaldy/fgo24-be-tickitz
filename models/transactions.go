package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"time"
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
