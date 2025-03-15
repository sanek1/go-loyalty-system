package entity

type User struct {
	ID       uint   `json:"ID"`
	Login    string `json:"Login"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Access   string
}
