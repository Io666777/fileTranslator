package repositories

import (
	"fileTranslator/cmd/models"
	"fileTranslator/cmd/storage"
)

func CreateUser(user models.User) (models.User, error) {
  db := storage.GetDB()
  sqlStatement := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`
  err := db.QueryRow(sqlStatement, user.Name, user.Password).Scan(&user.ID)
  if err != nil {
    return user, err
  }
  return user, nil
}

func 