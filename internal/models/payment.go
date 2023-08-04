package models

import (
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/database"
	"time"
)

type Payment struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	Value    float64   `json:"value"`
	Item     string    `json:"item"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Notes    string    `json:"notes"`
}

func GetPaymentsByUserID(userID int64) ([]Payment, error) {
	sqlStatement := `
		SELECT id, user_id, value, item, category, date, notes
		FROM payments
		WHERE user_id = $1
		ORDER BY date DESC`

	rows, err := database.DB.Query(sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		err = rows.Scan(&payment.ID, &payment.UserID, &payment.Value, &payment.Item, &payment.Category, &payment.Date, &payment.Notes)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return payments, nil
}

// InsertPayment inserts a new payment into the database
func InsertPayment(payment *Payment) error {
	// Save the payment to the database
	var id int64
	err := database.DB.QueryRow("INSERT INTO payments (value, item, category, date, notes, user_id) "+
		"VALUES ($1, $2, $3, $4, $5,$6) RETURNING id", payment.Value, payment.Item, payment.Category, payment.Date, payment.Notes, payment.UserID).Scan(&id)
	if err != nil {
		return err
	}
	payment.ID = id
	fmt.Println("New payment added with ID", payment.ID)
	return nil
}
