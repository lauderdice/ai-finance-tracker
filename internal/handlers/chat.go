package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/models"
	"github.com/lauderdice/ai-finance-tracker/internal/session"
	"html/template"
	"net/http"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	// Retrieve the currently logged-in user
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		// Get chat messages from the session
		fmt.Println("Fetching chat page for user", user.Email)
		messages, err := session.GetMessagesFromSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Display chat history
		tmpl := template.Must(template.ParseFiles("templates/chat.html"))
		tmpl.Execute(w, struct {
			Messages []models.ChatMessage
		}{
			Messages: messages,
		})
	} else if r.Method == "POST" {
		// Handle chat message
		r.ParseForm()
		message := r.Form.Get("message")
		message = "User: " + message
		// TODO: Generate bot response using OpenAI API
		botResponse := "Bot: " + "Random OpenAI response"

		// Add user message and bot response to session
		err = session.AddMessageToSession(w, r, models.ChatMessage{Sender: "User", Message: message})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = session.AddMessageToSession(w, r, models.ChatMessage{Sender: "User", Message: botResponse})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a struct to hold the response data
		responseData := struct {
			UserMessage string `json:"userMessage"`
			BotResponse string `json:"botResponse"`
		}{
			UserMessage: message,
			BotResponse: botResponse,
		}

		// Convert the response data to JSON and write it to the response
		responseBytes, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBytes)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Other code...
