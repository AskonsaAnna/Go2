package handlers

import (
	"database/sql"
	"lionforum/pkg/session" // replace with your actual path
	"log"
	"net/http"
)

var redirectURL string

func (app *DBRegister) Signin(w http.ResponseWriter, r *http.Request) {
	emailInput := r.FormValue("E-mail")
	passwordInput := r.FormValue("password")

	// Prepare data structure for template
	data := PageData{
		CurrentURL: r.URL.String(),
	}

	urlValue := r.FormValue("redirect")
	if urlValue != "/login" {
		redirectURL = urlValue
	}

	query := `SELECT id, username, email, password FROM users WHERE email =?`
	row := app.DB.QueryRow(query, emailInput)

	var userId int
	var username, email, password string

	err := row.Scan(&userId, &username, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			data.ErrorMessage = "Invalid email or password"
			renderTemplate(w, "register.html", data)
			return
		} else {
			data.ErrorMessage = "Internal server error"
			renderTemplate(w, "register.html", data)
			log.Println(err)
			return
		}
	}

	if CheckPassword(password, passwordInput) {
		sessionID := session.GenerateUUID()
		session.CreateNewSession(userId, sessionID, username, email)
		session.CreateSessionCookie(w, sessionID)
	} else {
		data.ErrorMessage = "Incorrect password"
		renderTemplate(w, "register.html", data)
		return
	}

	if redirectURL == "" {
		redirectURL = "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)

}
