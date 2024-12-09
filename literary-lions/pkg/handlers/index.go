package handlers

import (
	"lionforum/pkg/session"
	"net/http"
)

var globalData PageData

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the URL contains the parameter "showTopics=true"
	showTopics := r.URL.Query().Get("showTopics") == "true"

	sessionID, _ := session.ReadSessionCookie(r)
	_, activeuser, _ := session.GetActiveUser(sessionID)

	data := PageData{
		ShowTopics: showTopics,
		Username:   activeuser,
	}
	data.Topics = []string{"Fiction", "Crime/Thriller", "Science Fiction",
		"Fantasy", "Historical Fiction", "Romantic Fiction", "Poetry", "Non-fiction", "Other"}

	renderTemplate(w, "index.html", data)

}
