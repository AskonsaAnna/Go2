package session

import (
	"net/http"
)

func CreateSessionCookie(w http.ResponseWriter, uuid string) error {
	c := &http.Cookie{
		Name:     "session_id",
		Value:    uuid,
		MaxAge:   3600,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, c)

	return nil
}

func ReadSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}

	sessionID := cookie.Value

	return sessionID, nil
}

func DeleteSessionCookie(w http.ResponseWriter, sessionID string) {
	c := &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	sessionStore.Delete(sessionID)
	http.SetCookie(w, c)
}
