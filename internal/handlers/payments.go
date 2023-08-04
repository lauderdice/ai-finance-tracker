package handlers

import (
	"encoding/json"
	"github.com/lauderdice/ai-finance-tracker/internal/database"
	"github.com/lauderdice/ai-finance-tracker/internal/models"
	"github.com/lauderdice/ai-finance-tracker/internal/session"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func AddPayment(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if r.Method == "GET" {
		// Display the form to add a new payment
		t, err := template.ParseFiles("templates/add_payment.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		// Parse the form data and save the new payment
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		value, err := parseFloat(r.PostFormValue("value"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item := r.PostFormValue("item")
		category := r.PostFormValue("category")
		dateStr := r.PostFormValue("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		notes := r.PostFormValue("notes")
		payment := &models.Payment{
			Value:    value,
			Item:     item,
			Category: category,
			Date:     date,
			Notes:    notes,
			UserID:   user.ID,
		}
		// Save the payment to the database
		err = models.InsertPayment(payment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetPayments(w http.ResponseWriter, r *http.Request) {
	// Retrieve the currently logged-in user
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Retrieve all payments for the logged-in user from the database
	rows, err := database.DB.Query("SELECT id, value, item, category, date, notes FROM payments WHERE user_id = $1", user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	payments := make([]models.Payment, 0)
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.Value, &payment.Item, &payment.Category, &payment.Date, &payment.Notes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check the request type (JSON or HTML) using a query parameter or any other method you prefer.
	requestType := r.URL.Query().Get("type")

	if requestType == "json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payments)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/payments.html"))
		tmpl.Execute(w, struct {
			Payments []models.Payment
		}{
			Payments: payments,
		})
	}
}

func parseFloat(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0, err
	}
	return value, nil
}
