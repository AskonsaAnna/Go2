package handlers

import (
	"database/sql"
	"fmt"
	"lionforum/pkg/session"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (app *DBRegister) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the selected topic from the URL parameters
	selectedTopic := r.URL.Query().Get("topic")
	filter := r.URL.Query().Get("filter")

	sessionID, _ := session.ReadSessionCookie(r)
	activeuserID, activeuser, _ := session.GetActiveUser(sessionID)

	if selectedTopic != "all" && selectedTopic != "" {
		filter = "topic"
	}

	// Fetch all posts from the database
	allPosts := app.FetchPost()

	globalData.Posts = allPosts

	var filteredPosts []Post

	switch filter {
	case "topic":
		filteredPosts = FilterByTopic(selectedTopic)
	case "latest":
		filteredPosts = FilterByDate()
	case "views":
		filteredPosts, _ = app.FilterByViews()
	case "user":
		filteredPosts = FilterByUser(activeuserID)
	case "like":
		filteredPosts = app.FilterByLike(activeuserID)
	default:
		filteredPosts = allPosts
	}

	topics := []string{"Fiction", "Crime/Thriller", "Science Fiction",
		"Fantasy", "Historical Fiction", "Romantic Fiction", "Poetry", "Non-fiction", "Other"}

	showPost := r.URL.Query().Get("form") != "new"

	data := PageData{
		ShowPost:   showPost,
		Topics:     topics,
		Posts:      filteredPosts, // Use filtered posts here
		Username:   activeuser,
		CurrentURL: r.Header.Get("Host") + r.URL.Path + "?" + r.URL.RawQuery,
	}

	renderTemplate(w, "post.html", data)

}

func (app *DBRegister) NewPost(w http.ResponseWriter, r *http.Request) {

	// only allow registered users to add posts

	sessionID, _ := session.ReadSessionCookie(r)

	if !session.Authenticate(w, r, sessionID) {
		// http.Error(w, "Must be signed in to post", http.StatusUnauthorized)
		data := PageData{
			ShowPost:     false,
			ErrorMessage: "Must be signed in to post",
		}
		renderTemplate(w, "post.html", data)
		return
	}

	userID, username, _ := session.GetActiveUser(sessionID)

	posting := Post{
		UserID:  userID,
		Title:   r.FormValue("title"),
		Topic:   r.FormValue("topic"),
		Content: r.FormValue("posti"),
		User:    username,
	}

	trimmedContent := strings.TrimSpace(posting.Content)

	if len(trimmedContent) == 0 {
		data := PageData{
			ShowPost:     false,
			ErrorMessage: "Post content cannot be empty or just whitespaces",
		}
		renderTemplate(w, "post.html", data)
		return
	}

	_, err := app.DB.Exec("INSERT INTO posts (user_id, username, title, topic, content) VALUES (?, ?, ?, ?, ?)", posting.UserID, posting.User, posting.Title, posting.Topic, posting.Content)

	if err != nil {
		msg := fmt.Sprintf("Failed to add post to database: %v", err)
		data := PageData{
			ShowPost:     false,
			ErrorMessage: msg,
		}
		renderTemplate(w, "post.html", data)
		return
	}

	http.Redirect(w, r, "/topic", http.StatusSeeOther)

}

func (app *DBRegister) EditPost(w http.ResponseWriter, r *http.Request) {

	// log.Println("Edit post called")

	postID := r.FormValue("id")
	content := r.FormValue("content")

	trimmedContent := strings.TrimSpace(content)

	if len(trimmedContent) == 0 {
		http.Redirect(w, r, "/chat?id="+postID+"&error="+url.QueryEscape("Content cannot be empty or just whitespaces"), http.StatusSeeOther)
		return
	}

	stmt, err := app.DB.Prepare("UPDATE posts SET content = ?, edited_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, postID)
	if err != nil {
		msg := fmt.Sprintf("Failed to add post to database: %v", err)
		log.Println(msg)
		return
	}

	http.Redirect(w, r, "/chat?id="+postID, http.StatusSeeOther)
}

func (app *DBRegister) DeletePost(w http.ResponseWriter, r *http.Request) {

	// log.Println("DeletePost called")

	postID := r.FormValue("id")

	stmt, err := app.DB.Prepare("UPDATE posts SET title = 'Deleted post', content = '' WHERE id = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(postID)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	http.Redirect(w, r, "/chat?id="+postID, http.StatusSeeOther)

}

func (app *DBRegister) FetchPost() []Post {

	query := `
        SELECT id, user_id, title, topic, content, username, views, likes, dislikes, date
        FROM posts
    `

	rows, err := app.DB.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var posts []Post

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

		if post.Content == "" {
			post.Deleted = true
			post.Content = "Post has been deleted by the author"
		}

		post.Count, _ = app.CountCommentsByPostID(post.ID)

		localTime := post.Date.In(time.Local)

		// Форматируем дату в нужном формате
		formattedDate := localTime.Format("02.01.2006 at 15:04")
		post.DateString = formattedDate

		posts = append(posts, post)
	}

	return posts

}

// called, when chat.html is clicked
func (app *DBRegister) FetchPostByID(postID int) (*Post, error) {
	var post Post
	var editedAt sql.NullTime // Using sql.NullTime to handle NULL values in Go

	query := `
        SELECT id, user_id, title, topic, content, username, views, likes, dislikes, date, edited_at
        FROM posts
        WHERE id = ?
    `

	err := app.DB.QueryRow(query, postID).Scan(
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
		&editedAt,
	)
	if err != nil {
		log.Println("Error scanning row:", err)
		// Обработка ошибки
	}

	if post.Content == "" {
		post.Deleted = true
		post.Content = "Post has been deleted by the author"
	}

	// Применение временной зоны к дате
	localTime := post.Date.In(time.Local)

	// Форматирование даты с учетом локального времени
	formattedDate := localTime.Format("02.01.2006 at 15:04")
	post.DateString = formattedDate

	// If edited_at is not NULL, assign it to post.Edited
	if editedAt.Valid {
		post.Edited = editedAt.Time.In(time.Local)
		post.EditedString = post.Edited.Format("02.01.2006 at 15:04")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// No post found with the given ID
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}
		return nil, err
	}

	return &post, nil
}
