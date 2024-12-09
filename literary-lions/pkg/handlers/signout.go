package handlers

import (
	"lionforum/pkg/session"
	"net/http"
)

// add UUID
func (app *DBRegister) Signout(w http.ResponseWriter, r *http.Request) {

	sessionID, _ := session.ReadSessionCookie(r)
	session.DeleteSessionCookie(w, sessionID)

	http.Redirect(w, r, "/", http.StatusFound)
}
