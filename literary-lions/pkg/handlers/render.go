package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

// Helper function to render template with PageData, or another type of struct
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("web/template", filepath.Clean(tmpl))
	footer := "web/template/footer.html"

	t, err := template.ParseFiles(tmplPath, footer)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}
