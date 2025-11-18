package models

type Filestr struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	NAMEFILE string `json:"namefile"`
	AUTHOR   string `json:"author"`
	CLOUDKEY string `json:"cloudkey"`
	LANGFROM string `json:"langfrom"`
	LANGTO   string `json:"langto"`
}

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
