package handlers

import (
	"database/sql"
	"time"
)

type DBRegister struct {
	DB *sql.DB
}

type Post struct {
	ID           int
	UserID       int
	Title        string
	Topic        string
	Content      string
	User         string
	Views        int
	Comments     []Comment
	Likes        int
	Dislikes     int
	ActiveUser   string
	Date         time.Time
	Edited       time.Time
	Count        int
	DateString   string // Отформатированная дата
	EditedString string
	CurrentURL   string
	ErrorMessage string
	MyPost       bool
	Deleted      bool
}

type Comment struct {
	ID           int
	Content      string
	Date         time.Time
	Edited       time.Time
	DateString   string // Отформатированная дата
	EditedString string
	PostID       int
	CommentID    int // use for replies
	UserID       int
	Username     string
	Likes        int
	Dislikes     int
	//add replies
	Replies   []Comment
	PostTitle string
	MyPost    bool
	Deleted   bool
}

type PageData struct {
	ShowTopics       bool
	Results          []Comment
	ShowRegisterForm bool
	ShowPost         bool
	Topics           []string
	Posts            []Post
	Username         string
	CurrentURL       string
	Email            string
	ErrorMessage     string
	UsernameError    string // Ошибка для поля Username
	EmailError       string // Ошибка для поля Email
	PasswordError    string
}

type ProfileData struct {
	Username string
	Email    string

	EditEmail     bool
	EditPassword  bool
	EmailError    string
	PasswordError string
}
