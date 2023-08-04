package handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/models"
	"github.com/lauderdice/ai-finance-tracker/internal/session"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

func GetUploadScreen(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/upload_csv.html"))
	tmpl.Execute(w, nil)
}

func ProcessCSVUpload(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'   // Set the delimiter, default is ','
	reader.Comment = '#' // Set the comment character, lines starting with this character will be ignored
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		http.Error(w, "Error parsing CSV file", http.StatusInternalServerError)
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error parsing CSV file", http.StatusInternalServerError)
			return
		}

		value, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing value: %s", err.Error()), http.StatusBadRequest)
			return
		}
		item := record[0]
		category := record[1]
		date, err := time.Parse("2006-01-02", record[3])
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing date: %s", err.Error()), http.StatusBadRequest)
			return
		}
		notes := record[4]

		payment := &models.Payment{
			Value:    value,
			Item:     item,
			Category: category,
			Date:     date,
			Notes:    notes,
			UserID:   user.ID,
		}

		err = models.InsertPayment(payment)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error saving payment: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
