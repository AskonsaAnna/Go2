package main

import (
	"lionforum/pkg/database"
	"lionforum/pkg/handlers"
	"log"
	"net/http"
)

func main() {

	database, err := database.InitDB("./pkg/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./web/static/"))

	dbHandler := &handlers.DBRegister{
		DB: database,
	}

	// Register handler functions
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/search", dbHandler.SearchHandler)
	http.HandleFunc("/gotoregister", handlers.RegisterHandler)
	http.HandleFunc("/register", dbHandler.Register)
	http.HandleFunc("/login", dbHandler.Signin)
	http.HandleFunc("/logout", dbHandler.Signout)

	http.HandleFunc("/topic", dbHandler.PostHandler)
	http.HandleFunc("/newpost", dbHandler.NewPost)
	http.HandleFunc("/chat", dbHandler.ChatHandler)
	http.HandleFunc("/like", dbHandler.LikeHandler)
	http.HandleFunc("/dislike", dbHandler.DislikeHandler)
	http.HandleFunc("/post", dbHandler.PostHandler)
	http.HandleFunc("/views", dbHandler.ViewsHandler)
	http.HandleFunc("/addcomment", dbHandler.CommentHandler)
	http.HandleFunc("/reply", dbHandler.ReplyHandler)
	http.HandleFunc("/gotoprofile", handlers.ProfileHandler)
	http.HandleFunc("/update_email", dbHandler.UpdateEmail)
	http.HandleFunc("/update_password", dbHandler.UpdatePassword)
	http.HandleFunc("/edit_post", dbHandler.EditPost)
	http.HandleFunc("/delete_post", dbHandler.DeletePost)
	http.HandleFunc("/edit_comment", dbHandler.EditComment)
	http.HandleFunc("/delete_comment", dbHandler.DeleteComment)

	// Serve static files and templates
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the web server
	log.Println("Server started on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
