package main

type Todo struct {
	ID      uint `gorm:"primarykey"`
	Checked bool
	Text    string
}

type User struct {
	ID           uint `gorm:"primarykey"`
	Email        string
	PasswordHash string
}

type Session struct {
	ID     uint `gorm:"primarykey"`
	Email  string
	Cookie string
}
