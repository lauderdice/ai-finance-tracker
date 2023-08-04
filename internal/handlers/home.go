package handlers

import (
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/session"
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := session.GetUserFromSession(r)
	if err == nil {
		http.Redirect(w, r, "/get_payments", http.StatusFound)
		return
	}
	fmt.Printf("Fetching home page because user is not logged in\n")
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}
