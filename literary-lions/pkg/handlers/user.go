package handlers

import "database/sql"

func (app *DBRegister) FetchUsernameByID(userID int) string {
	var username string

	query := `
        SELECT username
        FROM users
        WHERE id = ?
    `

	err := app.DB.QueryRow(query, userID).Scan(
		&username,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// no user found with given id
			return "Deleted user"
		}

		return ""
	}

	return username
}
