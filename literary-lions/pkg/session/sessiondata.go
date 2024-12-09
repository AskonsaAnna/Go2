package session

import (
	"fmt"
	"time"
)

type Session struct {
	ID            string // Unique session ID (e.g., UUID)
	UserID        int    // ID of the user associated with this session
	Username      string
	Email         string
	CreatedAt     time.Time // Time when the session was created
	ExpiresAt     time.Time // Time when the session will expire
	EmailError    string
	PasswordError string
}

type SessionStore struct {
	// Example store, can be implemented using Redis, in-memory map, etc.
	store map[string]*Session
}

var sessionStore *SessionStore

func init() {
	sessionStore = &SessionStore{
		store: make(map[string]*Session),
	}
}

func CreateNewSession(userID int, uuid, username, email string) {
	newSession := Session{
		ID:        uuid,
		UserID:    userID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}

	sessionStore.Save(newSession.ID, &newSession)
}

// To interface with sessionStore from outside the session package
// give access to the ID, name and email of active user
func GetActiveUser(sessionID string) (int, string, string) {
	currentSession, err := sessionStore.Get(sessionID)

	if err != nil {
		return 0, "", ""
	}

	return currentSession.UserID, currentSession.Username, currentSession.Email
}

func UpdateShownEmail(sessionID, email string) {
	sessionStore.store[sessionID].Email = email
}

func EmailError(sessionID, errorMessage string) {
	sessionStore.store[sessionID].EmailError = errorMessage
}

func PasswdError(sessionID, errorMessage string) {
	sessionStore.store[sessionID].PasswordError = errorMessage
}

func GetErrors(sessionID string) (string, string) {
	currentSession, _ := sessionStore.Get(sessionID)

	return currentSession.EmailError, currentSession.PasswordError
}

func ClearErrors(sessionID string) {
	currentSession, _ := sessionStore.Get(sessionID)

	currentSession.EmailError = ""
	currentSession.PasswordError = ""
}

func (s *SessionStore) Get(sessionID string) (*Session, error) {
	session, exists := s.store[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (s *SessionStore) Save(sessionID string, session *Session) {
	s.store[sessionID] = session
}

func (s *SessionStore) Delete(sessionID string) {
	delete(s.store, sessionID)
}
