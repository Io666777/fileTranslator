package storage

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq" // драйвер PostgreSQL
)

var db *sql.DB

func InitDB() {
    // Загружаем .env
    if err := godotenv.Load(); err != nil {
        log.Fatal("Ошибка загрузки .env файла")
    }

    // Читаем переменные окружения
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    // Формируем DSN
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        dbHost, dbUser, dbPass, dbName, dbPort)

    // Подключаемся к БД
    var err error
    db, err = sql.Open("postgres", dsn)
    if err != nil {
        panic(err.Error())
    }

    // Проверяем соединение
    if err = db.Ping(); err != nil {
        panic(err.Error())
    }

    fmt.Println("БД успешно подключена")
}

func GetDB() *sql.DB {
    return db
}
