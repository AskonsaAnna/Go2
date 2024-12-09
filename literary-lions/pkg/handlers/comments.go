package handlers

import (
	"database/sql"
	"fmt"
	"lionforum/pkg/session"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (app *DBRegister) CommentHandler(w http.ResponseWriter, r *http.Request) {

	sessionID, _ := session.ReadSessionCookie(r)
	activeuserID, _, _ := session.GetActiveUser(sessionID)

	idStr := r.FormValue("post_id")
	postID, _ := strconv.Atoi(idStr)

	// only allow registered users to add comments
	if !session.Authenticate(w, r, sessionID) {
		// http.Error(w, "Must be signed in to comment", http.StatusUnauthorized)
		http.Redirect(w, r, "/chat?id="+idStr+"&error="+url.QueryEscape("Must be signed in to comment"), http.StatusSeeOther)
		return
	}

	add := Comment{
		UserID:    activeuserID,
		Content:   r.FormValue("comment"),
		PostID:    postID,
		PostTitle: r.FormValue("post_title"),
	}

	trimmedContent := strings.TrimSpace(add.Content)

	if len(trimmedContent) == 0 {
		http.Redirect(w, r, "/chat?id="+idStr+"&error="+url.QueryEscape("Content cannot be empty or just whitespaces"), http.StatusSeeOther)
		return
	}

	_, err := app.DB.Exec("INSERT INTO comments (post_id, user_id, content, post_title) VALUES (?, ?, ?, ?)", add.PostID, add.UserID, add.Content, add.PostTitle)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add comment to database: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chat?id="+idStr, http.StatusSeeOther)
}

func (app *DBRegister) ReplyHandler(w http.ResponseWriter, r *http.Request) {

	sessionID, _ := session.ReadSessionCookie(r)
	activeuserID, _, _ := session.GetActiveUser(sessionID)

	idStr := r.FormValue("comment-id")
	postID := r.FormValue("post-id")
	commentID, _ := strconv.Atoi(idStr)

	// only allow registered users to add comments
	if !session.Authenticate(w, r, sessionID) {
		// http.Error(w, "Must be signed in to comment", http.StatusUnauthorized)
		http.Redirect(w, r, "/chat?id="+postID+"&error="+url.QueryEscape("Must be signed in to reply"), http.StatusSeeOther)
		return
	}

	add := Comment{
		UserID:    activeuserID,
		Content:   r.FormValue("reply-content"),
		CommentID: commentID,
		PostTitle: r.FormValue("post_title"),
	}

	trimmedContent := strings.TrimSpace(add.Content)

	if len(trimmedContent) == 0 {
		http.Redirect(w, r, "/chat?id="+idStr+"&error="+url.QueryEscape("Content cannot be empty or just whitespaces"), http.StatusSeeOther)
		return
	}

	_, err := app.DB.Exec("INSERT INTO comments (post_id, comment_id, user_id, content, post_title) VALUES (?, ?, ?, ?, ?)", 0, add.CommentID, add.UserID, add.Content, add.PostTitle)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add reply to database: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chat?id="+postID, http.StatusSeeOther)
}

func (app *DBRegister) EditComment(w http.ResponseWriter, r *http.Request) {

	// log.Println("Edit comment called")

	postID := r.FormValue("id")
	commentID := r.FormValue("comment-id")
	content := r.FormValue("comment-content")

	trimmedContent := strings.TrimSpace(content)

	if len(trimmedContent) == 0 {
		http.Redirect(w, r, "/chat?id="+postID+"&error="+url.QueryEscape("Content cannot be empty or just whitespaces"), http.StatusSeeOther)
		return
	}

	stmt, err := app.DB.Prepare("UPDATE comments SET content = ?, edited_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, commentID)
	if err != nil {
		msg := fmt.Sprintf("Failed to add post to database: %v", err)
		log.Println(msg)
		return
	}

	http.Redirect(w, r, "/chat?id="+postID, http.StatusSeeOther)
}

func (app *DBRegister) DeleteComment(w http.ResponseWriter, r *http.Request) {

	// log.Println("DeleteComment called")

	postID := r.FormValue("id")
	commentID := r.FormValue("comment-id")

	stmt, err := app.DB.Prepare("UPDATE comments SET content = 'Deleted comment' WHERE id = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(commentID)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	http.Redirect(w, r, "/chat?id="+postID, http.StatusSeeOther)

}

func (app *DBRegister) FetchCommentByID(postID int) []Comment {

	query := `
        SELECT id, post_id, post_title, user_id, content, likes, dislikes, created_at, edited_at
        FROM comments
        WHERE post_id = ?
    `

	rows, _ := app.DB.Query(query, postID)
	defer rows.Close()
	var comments []Comment
	var editedAt sql.NullTime

	for rows.Next() {
		var comment Comment

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.PostTitle,
			&comment.UserID,
			&comment.Content,
			&comment.Likes,
			&comment.Dislikes,
			&comment.Date,
			&editedAt,
		)

		if comment.Content == "Deleted comment" {
			comment.Deleted = true
		}

		localTime := comment.Date.In(time.Local)

		// Форматирование даты с учетом локального времени
		formattedDate := localTime.Format("02.01.2006 at 15:04")
		comment.DateString = formattedDate

		if editedAt.Valid {
			comment.Edited = editedAt.Time.In(time.Local)
			comment.EditedString = comment.Edited.Format("02.01.2006 at 15:04")
		}

		comment.Username = app.FetchUsernameByID(comment.UserID)

		comment.Replies = app.FetchReplyByID(comment.ID)

		if err != nil {
			fmt.Println((err))
		}
		comments = append(comments, comment)
	}
	return comments

}

func (app *DBRegister) FetchReplyByID(commentID int) []Comment {

	query := `
        SELECT id, user_id, content, created_at
        FROM comments
        WHERE comment_id = ?
    `

	rows, _ := app.DB.Query(query, commentID)
	defer rows.Close()
	var replies []Comment

	for rows.Next() {
		var reply Comment

		err := rows.Scan(
			&reply.ID,
			&reply.UserID,
			&reply.Content,
			&reply.Date,
		)

		localTime := reply.Date.In(time.Local)

		// Форматирование даты с учетом локального времени
		formattedDate := localTime.Format("02.01.2006 at 15:04")
		reply.DateString = formattedDate

		reply.Username = app.FetchUsernameByID(reply.UserID)

		// implement a fetch for replies
		// comment.Replies = append(comment.Replies, Comment{
		// 	Username: "Lion cub",
		// 	Content:  "test reply",
		// })

		if err != nil {
			fmt.Println((err))
		}
		replies = append(replies, reply)
	}
	return replies

}

func (app *DBRegister) CountCommentsByPostID(postID int) (int, error) {
	var count int

	query := `
        SELECT COUNT(*)
        FROM comments
        WHERE post_id = ?
    `

	// Выполнение запроса и получение результата
	err := app.DB.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Пример использования функции в обработчике
func (app *DBRegister) CommentCountHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID поста из запроса
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Подсчитываем количество комментариев
	count, err := app.CountCommentsByPostID(postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to count comments: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем количество комментариев как текстовый ответ
	fmt.Fprintf(w, "Total comments for post %d: %d", postID, count)
}
