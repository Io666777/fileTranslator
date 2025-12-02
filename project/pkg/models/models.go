package models

type User struct {
	ID       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required" db:"name"`
	Username string `json:"username" binding:"required" db:"username"`
	Password string `json:"password" binding:"required" db:"password_hash"`
}

type File struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" binding:"required"`
	Path        string `json:"path" db:"path"`
	Status      string `json:"status" db:"status"` 
	UserID      int    `json:"user_id" db:"user_id"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	FileContent []byte `json:"-" db:"file_content"`
}
