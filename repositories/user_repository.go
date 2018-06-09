package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"rest-api/models"
	"rest-api/utils/crypto"
)

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	const query = `
		select
			id,
			email,
			name
		from
			users
		where
			id = ?
	`
	var user models.User
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name)
	return &user, err
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	const query = `
		select
			id,
			email,
			name
		from
			users
		where
			email = ?
	`
	var user models.User
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name)
	return &user, err
}

func GetPrivateUserDetailsByEmail(db *sql.DB, email string) (*models.PrivateUserDetails, error) {
	const query = `
		select
			id,
			password,
			salt
		from
			users
		where
			email = ?
	`
	var u models.PrivateUserDetails
	err := db.QueryRow(query, email).Scan(&u.ID, &u.Password, &u.Salt)
	return &u, err
}

func CreateUser(db *sql.DB, email, name, password string) (int64, error) {
	const query = `
		insert into users (
			email,
			name,
			password,
			salt
		) values (
			?,
			?,
			?,
			?
		)
	`
	salt := crypto.GenerateSalt()
	hashedPassword := crypto.HashPassword(password, salt)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	res, err := stmt.Exec(email, name, hashedPassword, salt)
	if err != nil {
		log.Fatal("Cannot run insert statement", err)
	}
	id, _ := res.LastInsertId()
	fmt.Printf("Inserted row: %d", id)

	return id, err
}

func GetUsers(db *sql.DB, page, resultsPerPage int) ([]*models.User, error) {
	const query = `
		select
			id,
			name,
			email
		from
			users
		limit ? offset ?
		`
	users := make([]*models.User, 0)
	offset := (page - 1) * resultsPerPage

	rows, err := db.Query(query, resultsPerPage, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Email, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, err
}
