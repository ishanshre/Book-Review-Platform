package models

import "time"

type User struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Gender            string    `json:"gender"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
	ProfilePic        string    `json:"profile_pic"`
	CitizenshipNumber string    `json:"citizenship_number"`
	CitizenshipFront  string    `json:"citizenship_front"`
	CitizenshipBack   string    `json:"citizenship_back"`
	AccessLevel       int       `json:"access_level"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastLogin         time.Time `json:"last_login"`
}
