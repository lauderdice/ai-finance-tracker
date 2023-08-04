package handlers

import (
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/models"
	"github.com/lauderdice/ai-finance-tracker/internal/session"
	"html/template"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint handler Register was called")
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {

		email := r.FormValue("email")
		password := r.FormValue("password")

		err := models.CreateUser(email, password)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		fmt.Printf("User with email %s was created\n", email)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint handler Login was called")
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := models.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = models.VerifyPassword(user.Password, password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		fmt.Printf("User with email %s logged in\n", email)
		// Set a secure cookie
		err = session.SetUserSession(w, r, user.ID)
		if err != nil {
			fmt.Println("Error setting user session:", err)
			http.Error(w, "Unable to set session", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := session.ClearUserSession(w, r)
	if err != nil {
		fmt.Println("Error clearing user session:", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
