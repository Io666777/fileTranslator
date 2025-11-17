package models

type Filestr struct {
    ID       string `gorm:"primaryKey" json:"id"`
    NAMEFILE string `json:"namefile"`
    AUTHOR   string `json:"author"`
    CLOUDKEY string `json:"cloudkey"`
    LANGFROM string `json:"langfrom"`
    LANGTO   string `json:"langto"`
}

type User struct {
    ID       string `gorm:"primaryKey" json:"id"`
    Name     string `json:"name"`
    Email    string `gorm:"unique" json:"email"`
    Password string `json:"password"`
}
