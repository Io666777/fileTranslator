package filetranslator

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Filebt struct {
	Id     int    `json:"id"`
	Titleb string `json:"title"`
}

type Fileat struct {
	Id     int    `json:"id"`
	Titlea string `json:"title"`
}
