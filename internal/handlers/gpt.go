// handlers.go
package handlers

import (
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/models"
	"github.com/lauderdice/ai-finance-tracker/internal/session"
	"net/http"

	"strings"
)

var apiKey = ""

func LLMInitialize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	payments, err := models.GetPaymentsByUserID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contextString := createTableContext(payments)
	response, err := callOpenAIAPI(contextString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func createTableContext(payments []models.Payment) string {
	var sb strings.Builder
	sb.WriteString("Payments Table:\n\n")
	sb.WriteString("Value\tItem\tCategory\tDate\tNotes\n")

	for _, payment := range payments {
		sb.WriteString(fmt.Sprintf("%.2f\t%s\t%s\t%s\t%s\n", payment.Value, payment.Item, payment.Category, payment.Date.Format("2006-01-02"), payment.Notes))
	}

	return sb.String()
}

func callOpenAIAPI(context string) ([]byte, error) {
	// Replace this with your actual OpenAI API call implementation
	// You can use the official OpenAI API client or any other preferred method
	// Pass the context string as input to the API
	return []byte("This is a mock response from the OpenAI API"), nil
}
