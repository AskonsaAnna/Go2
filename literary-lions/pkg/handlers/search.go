package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// SearchHandler handles the search requests
func (app *DBRegister) SearchHandler(w http.ResponseWriter, r *http.Request) {

	searchQuery := r.FormValue("query")
	if searchQuery == "" {
		http.Error(w, "Search query is empty", http.StatusBadRequest)
		return
	}

	// Использование базы данных для поиска постов по содержимому
	postQuery := `
	SELECT id, user_id, title, topic, content, username, views, likes, dislikes, date
	FROM posts
	WHERE content LIKE ?
	`

	rows, err := app.DB.Query(postQuery, "%"+strings.TrimSpace(searchQuery)+"%")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	var results []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Topic,
			&post.Content,
			&post.User,
			&post.Views,
			&post.Likes,
			&post.Dislikes,
			&post.Date,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		post.Count, _ = app.CountCommentsByPostID(post.ID)
		localTime := post.Date.In(time.Local)
		post.DateString = localTime.Format("02.01.2006 at 15:04")
		results = append(results, post)
	}

	commentQuery := `
	SELECT post_id, user_id, content, created_at, post_title
	FROM comments
	WHERE content LIKE ?
	`

	rows, err = app.DB.Query(commentQuery, "%"+strings.TrimSpace(searchQuery)+"%")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.Date,
			&comment.PostTitle,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		comment.Username = app.FetchUsernameByID(comment.UserID)

		localTime := comment.Date.In(time.Local)
		comment.DateString = localTime.Format("02.01.2006 at 15:04")
		comments = append(comments, comment)
	}

	data := PageData{
		ShowTopics: false, // Hide topics list in the search results page
		Posts:      results,
		ShowPost:   true,
		Results:    comments,
	}

	// Load and initialize the template for search results
	tmpl, err := template.ParseFiles(filepath.Join("web/template", "post.html"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}
