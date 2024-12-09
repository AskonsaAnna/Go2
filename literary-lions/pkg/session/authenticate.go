package session

import (
	"net/http"
)

// might need some changes
func Authenticate(w http.ResponseWriter, r *http.Request, uuid string) bool {
	sessionID, _ := ReadSessionCookie(r)

	currentSession, err := sessionStore.Get(sessionID)

	if err != nil {
		return false
	}

	if currentSession.ID == "" || sessionID != currentSession.ID {
		return false
	}

	return true
}
