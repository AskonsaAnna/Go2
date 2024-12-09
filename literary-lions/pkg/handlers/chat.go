package handlers

import (
	"lionforum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

func (app *DBRegister) ChatHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil || postID < 0 {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	sessionID, _ := session.ReadSessionCookie(r)
	activeID, activeuser, _ := session.GetActiveUser(sessionID)

	post, err := app.FetchPostByID(postID)

	if post.UserID == activeID || post.User == activeuser {
		post.MyPost = true
	}

	post.ErrorMessage = r.URL.Query().Get("error")

	if err != nil {
		// http.Error(w, "Post not found", http.StatusNotFound)
		log.Println(err)
	}

	post.Comments = app.FetchCommentByID(post.ID)

	for i, comment := range post.Comments {
		if comment.Username == activeuser {
			post.Comments[i].MyPost = true
		}
	}

	post.CurrentURL = r.Header.Get("Host") + r.URL.Path + "?" + r.URL.RawQuery

	post.ActiveUser = activeuser

	renderTemplate(w, "chat.html", post)

}

func (app *DBRegister) ViewsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, _ := strconv.Atoi(idStr)

	stmt, err := app.DB.Prepare("UPDATE posts SET views = views + 1 WHERE id = ?")
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
	// Перенаправление на страницу с обновленными счетчиками
	http.Redirect(w, r, "/chat?id="+idStr, http.StatusSeeOther)
}
