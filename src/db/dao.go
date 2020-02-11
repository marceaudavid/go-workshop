package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/marceaudavid/learn-go/src/models"
	uuid "github.com/satori/go.uuid"
)

// Connect ...
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return db, nil
}

// Reset ...
func Reset() {
	db, _ := Connect()
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS tokens")
	db.Exec("CREATE TABLE users (id SERIAL,username VARCHAR(255), password VARCHAR(255), email VARCHAR(255), data JSON, PRIMARY KEY (id))")
	db.Exec("CREATE TABLE tokens (id VARCHAR(255), user_id INT, expiration DATETIME)")
	defer db.Close()
}

// GetUser ...
func GetUser(body models.User) (models.User, error) {
	db, _ := Connect()
	defer db.Close()
	var user models.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE users.email = ?", body.Email).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		fmt.Printf(err.Error())
		return user, err
	}
	return user, nil
}

// CreateUser ...
func CreateUser(user models.User, hash string) (sql.Result, error) {
	db, _ := Connect()
	defer db.Close()
	res, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", user.Username, hash, user.Email)
	if err != nil {
		return nil, err
	}
	return res, err
}

// GetData ...
func GetData(id int) (string, error) {
	db, _ := Connect()
	defer db.Close()
	var json string = "{}"
	db.QueryRow("SELECT data FROM users WHERE id = ?", id).Scan(&json)
	return json, nil
}

// SaveData ...
func SaveData(json string, id int) (sql.Result, error) {
	db, _ := Connect()
	defer db.Close()
	res, err := db.Exec("UPDATE users SET data = ? WHERE id = ?", json, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UserExist ...
func UserExist(user models.User) (bool, error) {
	db, _ := Connect()
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", user.Username).Scan(&username)
	if err != nil {
		return true, err
	}
	return false, nil
}

// CreateToken ...
func CreateToken(user models.User, expiration int64) (*string, error) {
	db, _ := Connect()
	defer db.Close()
	token := uuid.NewV4().String()
	date := time.Now().Add(time.Duration(expiration) * time.Minute)
	_, err := db.Exec("INSERT INTO tokens (id, user_id, expiration) VALUES (?, ?, ?)", token, user.ID, date)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// FromToken ...
func FromToken(token string) models.User {
	db, _ := Connect()
	defer db.Close()
	var user models.User
	db.QueryRow("SELECT users.id, users.username, users.password, users.email FROM users JOIN tokens ON users.id = tokens.user_id WHERE tokens.id = ?", token).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user
}

// DeleteTokens ...
func DeleteTokens(cookie *http.Cookie) (sql.Result, error) {
	db, _ := Connect()
	defer db.Close()
	res, err := db.Exec("DELETE FROM tokens WHERE id = ?", cookie.Value)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteExpiredTokens ...
func DeleteExpiredTokens() (*int64, error) {
	db, _ := Connect()
	defer db.Close()
	res, err := db.Exec("DELETE FROM tokens WHERE expiration < NOW()")
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
