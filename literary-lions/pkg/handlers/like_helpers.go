package handlers

import "log"

func (app *DBRegister) addLike(id, user_id int) {
	// Prepare the SQL statement for incrementing the likes
	stmt, err := app.DB.Prepare("UPDATE posts SET likes = likes + 1 WHERE id = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(id)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	stmt, err = app.DB.Prepare("INSERT INTO post_likes (post_id, user_id, like_type) VALUES (?, ?, ?)")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, user_id, "like")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}
}

func (app *DBRegister) addDislike(id, user_id int) {

	stmt, err := app.DB.Prepare("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	stmt, err = app.DB.Prepare("INSERT INTO post_likes (post_id, user_id, like_type) VALUES (?, ?, ?)")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, user_id, "dislike")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

}

func (app *DBRegister) removeLike(id, user_id int) {
	// Prepare the SQL statement for decrementing the likes
	stmt, err := app.DB.Prepare("UPDATE posts SET likes = likes - 1 WHERE id = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(id)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	stmt, err = app.DB.Prepare("DELETE FROM post_likes WHERE post_id = ? AND user_id = ? AND like_type = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, user_id, "like")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}
}

func (app *DBRegister) removeDislike(id, user_id int) {

	stmt, err := app.DB.Prepare("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

	stmt, err = app.DB.Prepare("DELETE FROM post_likes WHERE post_id = ? AND user_id = ? AND like_type = 'dislike'")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, user_id)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing SQL statement:", err)
		return
	}

}
