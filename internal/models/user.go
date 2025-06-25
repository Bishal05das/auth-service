package models

import "time"

type User struct {
	ID       int
	Email    string
    PasswordHash  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}


type UserRegister struct {
	Email     string
	Password  string
}

type UserLogin struct {
	Email    string
	Password string
}

