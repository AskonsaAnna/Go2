package handlers

import (
	"database/sql"
	"fmt"
	"lionforum/pkg/session"
	"log"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	emailBoolStr := r.FormValue("edit_email")
	var emailBool bool

	if emailBoolStr == "true" {
		emailBool = true
	}

	passwordBoolStr := r.FormValue("edit_password")
	var passwordBool bool

	if passwordBoolStr == "true" {
		passwordBool = true
	}
	sessionID, _ := session.ReadSessionCookie(r)
	_, activeuser, activeemail := session.GetActiveUser(sessionID)

	data := ProfileData{
		Username:     activeuser,
		Email:        activeemail,
		EditEmail:    emailBool,
		EditPassword: passwordBool,
	}

	// errEmail, errPasswd := session.GetErrors(sessionID)

	data.EmailError, data.PasswordError = session.GetErrors(sessionID)

	if data.EmailError != "" {
		data.EditEmail = true
	}

	if data.PasswordError != "" {
		data.EditPassword = true
	}

	session.ClearErrors(sessionID)

	renderTemplate(w, "profile.html", data)
}

func (app *DBRegister) UpdateEmail(w http.ResponseWriter, r *http.Request) {

	editEmail := r.FormValue("email")

	sessionID, _ := session.ReadSessionCookie(r)
	activeuserID, _, _ := session.GetActiveUser(sessionID)

	err := app.DB.QueryRow("SELECT email FROM users WHERE email = ?", editEmail).Scan(&editEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			// Пользователь не найден, можно продолжить регистрацию
			// Здесь можно добавить код для создания нового пользователя
		} else {
			// Произошла другая ошибка при выполнении запроса
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	} else {
		// Пользователь с таким email уже существует
		// http.Error(w, "Email already in use", http.StatusBadRequest)
		session.EmailError(sessionID, "Email already in use")
		http.Redirect(w, r, "/gotoprofile", http.StatusSeeOther)
		return
	}

	stmt, err := app.DB.Prepare("UPDATE users SET email = ? WHERE id = ?")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(editEmail, activeuserID)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	session.UpdateShownEmail(sessionID, editEmail)
	log.Println("email updated")
	http.Redirect(w, r, "/gotoprofile", http.StatusSeeOther)
}

func (app *DBRegister) UpdatePassword(w http.ResponseWriter, r *http.Request) {

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	sessionID, _ := session.ReadSessionCookie(r)
	activeuserID, _, _ := session.GetActiveUser(sessionID)

	query := `
        SELECT password
        FROM users
		WHERE id = ?
    `
	var checkPassword string
	err := app.DB.QueryRow(query, activeuserID).Scan(
		&checkPassword,
	)
	if err != nil {
		fmt.Println("Error fetching from database")
		return
	}
	if CheckPassword(checkPassword, currentPassword) {

		stmt, err := app.DB.Prepare("UPDATE users SET password = ? WHERE id = ?")
		if err != nil {
			log.Println("Error preparing SQL statement:", err)
			return
		}
		defer stmt.Close()

		if newPassword == confirmPassword {

			hashedPassword, _ := HashPassword(newPassword)

			_, err = stmt.Exec(hashedPassword, activeuserID)
			if err != nil {
				// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error executing SQL statement:", err)
				return
			}

			log.Println("password updated")

		} else {
			session.PasswdError(sessionID, "Passwords do not match")
		}

		// http.Redirect(w, r, "/gotoprofile", http.StatusSeeOther)
	} else {
		// http.Error(w, "incorrect password", http.StatusUnauthorized)
		session.PasswdError(sessionID, "Incorrect password")
		// http.Redirect(w, r, "/gotoprofile", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/gotoprofile", http.StatusSeeOther)

}
