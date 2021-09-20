package main

type Todo struct {
	ID      uint `gorm:"primarykey"`
	Checked bool
	Text    string
}
