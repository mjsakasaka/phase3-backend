package models

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"size:100;unique"`
	Password string `json:"password" gorm:"size:100"`
	Email    string `json:"email" gorm:"size:100"`
}
