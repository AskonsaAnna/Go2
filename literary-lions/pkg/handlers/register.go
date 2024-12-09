package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	showRegisterForm := r.URL.Query().Get("form") == "register"

	data := PageData{
		ShowRegisterForm: showRegisterForm,
		CurrentURL:       r.FormValue("redirect"),
	}

	tmpl, err := template.ParseFiles(filepath.Join("web/template", "register.html"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}

func renderTemplateWithError(w http.ResponseWriter, tmpl string, data PageData, errorMessage string) {
	tmplPath := filepath.Join("web", "template", filepath.Clean(tmpl))

	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	// Если есть ошибка, передаем её в data и добавляем JavaScript для alert()
	if errorMessage != "" {
		data.ErrorMessage = errorMessage
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}

func (app *DBRegister) Register(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	email := r.FormValue("email")
	confirmEmail := r.FormValue("confirm-email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	var data PageData

	data.ShowRegisterForm = true

	var newUserName string
	err := app.DB.QueryRow("SELECT username FROM users WHERE username = ?", userName).Scan(&newUserName)
	if err != nil && err != sql.ErrNoRows {
		renderTemplateWithError(w, "register.html", data, "Server error. Please try again later.")
		return
	}

	if newUserName != "" {
		renderTemplateWithError(w, "register.html", data, "Username already in use")
		return
	}

	if email != confirmEmail {
		renderTemplateWithError(w, "register.html", data, "Emails do not match")
		return
	}

	var newUserEmail string
	err = app.DB.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&newUserEmail)
	if err != nil && err != sql.ErrNoRows {
		renderTemplateWithError(w, "register.html", data, "Server error. Please try again later.")
		return
	}

	if newUserEmail != "" {
		renderTemplateWithError(w, "register.html", data, "An account with this email already exists")
		return
	}

	if password != confirmPassword {
		renderTemplateWithError(w, "register.html", data, "Passwords do not match")
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		renderTemplateWithError(w, "register.html", data, "Failed to hash password")
		return
	}

	_, err = app.DB.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", userName, hashedPassword, email)
	if err != nil {
		renderTemplateWithError(w, "register.html", data, fmt.Sprintf("Failed to register user: %v", err))
		return
	}

	data.ShowRegisterForm = false

	redirectURL := fmt.Sprintf("/login?E-mail=%s&password=%s", email, password)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
