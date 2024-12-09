package handlers

import (
	"database/sql"
	"lionforum/pkg/session"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (app *DBRegister) LikeHandler(w http.ResponseWriter, r *http.Request) {

	// if post_id exists, then handle post_likes
	// else (comment_id) handle comment_likes

	// Get the post ID from the form data
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)

	sessionID, _ := session.ReadSessionCookie(r)

	// only allow registered users to add likes
	if !session.Authenticate(w, r, sessionID) {
		// http.Error(w, "Must be signed in to like", http.StatusUnauthorized)
		http.Redirect(w, r, "/chat?id="+idStr+"&error="+url.QueryEscape("Must be signed in to like"), http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {

		commentIDstr := r.FormValue("comment-id")
		commentID, _ := strconv.Atoi(commentIDstr)

		// fmt.Println(id, commentID)

		if err != nil {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}

		var likeType string

		activeuserid, _, _ := session.GetActiveUser(sessionID)

		if commentID == 0 {
			err = app.DB.QueryRow("SELECT like_type FROM post_likes WHERE post_id = ? AND user_id = ?", id, activeuserid).Scan(&likeType)

			if err != nil {

				// if no like or dislike
				if err == sql.ErrNoRows {
					app.addLike(id, activeuserid)
				} else {
					// Handle other errors
					log.Println("Error querying like_type:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {

				// already liked
				if likeType == "like" {
					app.removeLike(id, activeuserid)

					// already disliked
				} else {
					app.removeDislike(id, activeuserid)
					app.addLike(id, activeuserid)
				}
			}

		} else {
			err = app.DB.QueryRow("SELECT like_type FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, activeuserid).Scan(&likeType)

			if err != nil {

				// if no like or dislike
				if err == sql.ErrNoRows {
					app.addCommentLike(commentID, activeuserid)
				} else {
					// Handle other errors
					log.Println("Error querying like_type:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {

				// already liked
				if likeType == "like" {
					app.removeCommentLike(commentID, activeuserid)

					// already disliked
				} else {
					app.removeCommentDislike(commentID, activeuserid)
					app.addCommentLike(commentID, activeuserid)
				}
			}
		}

		// Redirect to the page with updated counts
		http.Redirect(w, r, "/chat?id="+idStr, http.StatusSeeOther)
	}
}

func (app *DBRegister) DislikeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	commentIDstr := r.FormValue("comment-id")

	sessionID, _ := session.ReadSessionCookie(r)

	// only allow registered users to add dislikes
	if !session.Authenticate(w, r, sessionID) {
		// http.Error(w, "Must be signed in to dislike", http.StatusUnauthorized)
		http.Redirect(w, r, "/chat?id="+idStr+"&error="+url.QueryEscape("Must be signed in to dislike"), http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Получаем ID поста из данных формы

		id, err := strconv.Atoi(idStr)
		commentID, _ := strconv.Atoi(commentIDstr)
		if err != nil {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}
		var likeType string
		sessionID, _ := session.ReadSessionCookie(r)
		activeuserid, _, _ := session.GetActiveUser(sessionID)
		if commentID == 0 {

			err = app.DB.QueryRow("SELECT like_type FROM post_likes WHERE post_id = ? AND user_id = ?", id, activeuserid).Scan(&likeType)

			if err != nil {
				if err == sql.ErrNoRows {
					app.addDislike(id, activeuserid)
				} else {
					// Handle other errors
					log.Println("Error querying like_type:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {

				// Handle likeType
				if likeType == "dislike" {
					app.removeDislike(id, activeuserid)
				} else {
					app.removeLike(id, activeuserid)
					app.addDislike(id, activeuserid)
				}
			}

		} else {
			err = app.DB.QueryRow("SELECT like_type FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, activeuserid).Scan(&likeType)

			if err != nil {
				if err == sql.ErrNoRows {
					app.addCommentDislike(commentID, activeuserid)
				} else {
					// Handle other errors
					log.Println("Error querying like_type:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {

				// Handle likeType
				if likeType == "dislike" {
					app.removeCommentDislike(commentID, activeuserid)
				} else {
					app.removeCommentLike(commentID, activeuserid)
					app.addCommentDislike(commentID, activeuserid)
				}
			}

		}

		// Перенаправление на страницу с обновленными счетчиками
		http.Redirect(w, r, "/chat?id="+idStr, http.StatusSeeOther)
	}

}
