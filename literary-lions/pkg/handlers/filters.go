package handlers

import (
	"log"
	"time"
)

func FilterByTopic(topic string) []Post {
	var filteredPosts []Post

	allPosts := globalData.Posts

	for _, post := range allPosts {
		if post.Topic == topic {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

func FilterByDate() []Post {
	var filteredPosts []Post

	allPosts := globalData.Posts

	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	for _, post := range allPosts {
		if post.Date.After(oneWeekAgo) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

func (app *DBRegister) FilterByViews() ([]Post, error) {
	// SQL-запрос для выборки двух постов с наибольшим количеством просмотров
	query := `
        SELECT id, title, topic, content, views, date, likes, dislikes, username
        FROM posts
        ORDER BY views DESC
        LIMIT 2;
    `

	// Выполнение запроса
	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Слайс для хранения результатов
	var topPosts []Post

	// Итерация по результатам запроса
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Topic, &post.Content, &post.Views, &post.Date, &post.Likes, &post.Dislikes, &post.User); err != nil {
			return nil, err
		}

		post.Count, _ = app.CountCommentsByPostID(post.ID)
		topPosts = append(topPosts, post)
	}

	// Проверка на наличие ошибок после итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topPosts, nil
}

func FilterByUser(activeuserID int) []Post {
	var filteredPosts []Post

	allPosts := globalData.Posts

	for _, post := range allPosts {
		if post.UserID == activeuserID {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

func (app *DBRegister) FilterByLike(activeuserID int) []Post {
	var likedPosts []Post

	// Prepare the SQL query to select posts that the current user has liked
	query := `
	 SELECT p.id, p.title, p.topic, p.content, p.user_id, p.username, p.views, p.likes, p.dislikes, p.date
	 FROM posts p
	 INNER JOIN post_likes pl ON p.id = pl.post_id
	 WHERE pl.user_id = ? AND pl.like_type = 'like'
	`

	// Execute the query with the current user's ID
	rows, err := app.DB.Query(query, activeuserID)
	if err != nil {
		log.Println("Error fetching liked posts:", err)
		return likedPosts
	}
	defer rows.Close()

	// Iterate over the result set and populate the likedPosts slice
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Topic, &post.Content, &post.UserID, &post.User, &post.Views, &post.Likes, &post.Dislikes, &post.Date)
		if err != nil {
			log.Println("Error scanning post row:", err)
			continue
		}

		localTime := post.Date.In(time.Local)

		post.Count, _ = app.CountCommentsByPostID(post.ID)

		// Форматируем дату в нужном формате
		formattedDate := localTime.Format("02.01.2006 at 15:04")
		post.DateString = formattedDate

		likedPosts = append(likedPosts, post)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		log.Println("Error iterating over liked posts:", err)
	}

	return likedPosts
}
