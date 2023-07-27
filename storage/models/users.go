package models

import (
	"auth_service/storage"
	"auth_service/utils"
	"errors"
	"log"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (u *User) StoreToDB() error {

	if u.Username == "" || u.Email == "" || u.Password == "" {
		return errors.New("missing fields")
	}

	hashedPassword := utils.HashPassword(u.Password)
	uuidF := uuid.New()
	insert, err := storage.DB.Query("INSERT INTO Users(user_id, username, email, password) VALUES ( ?, ?, ?, ?)", uuidF, u.Username, u.Email, hashedPassword)

	if err != nil {
		return err
	}

	defer insert.Close()

	return nil
}

func GetUserByUsername(username string) (User, error) {

	result := storage.DB.QueryRow("SELECT * FROM Users WHERE username=?", username)

	resUser := User{}

	err := result.Scan(&resUser.Id, &resUser.Username, &resUser.Email, &resUser.Password)

	if err != nil {
		log.Println(err)
		return resUser, err
	}

	return resUser, nil
}
