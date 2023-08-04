package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/lauderdice/ai-finance-tracker/internal/models"
	"net/http"
)

var sessionStore = sessions.NewCookieStore([]byte("secret-key"))

func SetUserSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	// Create a new session for the user
	session, err := sessionStore.New(r, "user-session")
	if err != nil {
		return err
	}

	// Set the user ID as a session value
	session.Values["user_id"] = userID

	// Save the session to the cookie
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func GetUserSession(r *http.Request) (*sessions.Session, error) {
	// Retrieve the session for the current request
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		return nil, err
	}

	return session, nil
}

func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	// Retrieve the session for the current request
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		return err
	}

	// Clear the user ID from the session
	delete(session.Values, "user_id")
	delete(session.Values, "messages")
	// Save the updated session to the cookie
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func GetUserFromSession(r *http.Request) (*models.User, error) {
	// Retrieve the current session for the user
	userSession, err := GetUserSession(r)
	if err != nil {
		return nil, err
	}

	// Retrieve the user ID from the session
	userID, ok := userSession.Values["user_id"].(int64)
	if !ok {
		return nil, fmt.Errorf("user ID not found in session")
	}

	// Retrieve the user with the matching ID from the database
	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func AddMessageToSession(w http.ResponseWriter, r *http.Request, message models.ChatMessage) error {
	// Retrieve the session for the current request
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		return err
	}

	// Add the message to the session
	fmt.Println("Adding message to session: ", message)
	messages := session.Values["messages"]
	if messages == nil {
		messages = []models.ChatMessage{}
	}
	session.Values["messages"] = append(messages.([]models.ChatMessage), message)

	// Save the session to the cookie
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func GetMessagesFromSession(r *http.Request) ([]models.ChatMessage, error) {
	// Retrieve the session for the current request
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		return nil, err
	}

	// Retrieve the messages from the session
	fmt.Println("Retrieving messages from session: ", session.Values["messages"])
	messages := session.Values["messages"]
	if messages == nil {
		return []models.ChatMessage{}, nil
	}

	return messages.([]models.ChatMessage), nil
}
