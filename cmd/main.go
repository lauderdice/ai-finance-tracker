package main

import (
	"fmt"
	"github.com/lauderdice/ai-finance-tracker/internal/database"
	"github.com/lauderdice/ai-finance-tracker/internal/handlers"
	"log"
	"net/http"
)

func main() {
	err := database.InitDatabase()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	fmt.Println("Database initialized successfully")
	http.HandleFunc("/llm_initialize", handlers.LLMInitialize)
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/chat", handlers.Chat)

	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/add_payment", handlers.AddPayment)
	http.HandleFunc("/get_payments", handlers.GetPayments)
	http.HandleFunc("/upload_csv_form", handlers.GetUploadScreen)
	http.HandleFunc("/upload_csv", handlers.ProcessCSVUpload)
	http.HandleFunc("/logout", handlers.Logout)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
