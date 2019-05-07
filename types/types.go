package types

import "github.com/dgrijalva/jwt-go"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Id string `json:"id"`
}

type Artist struct {
	Name string `json:"name"`
	Id int `json:"id"`
}

type Set struct {
	User int `json:"user_id"`
	Location int `json:"location_id"`
	Artist int `json:"artist_id"`
	Id int `json:"id"`
}

type Location struct {
	Name string `json:"name"`
	Description string `json:"description"`
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
	Id int `json:"id"`
}

type Claims struct {
	Username string `json:"username"`
	Role string `json:"role"`
	Id string `json:"id"`
	jwt.StandardClaims
}

type Error interface {
	StatusCode() int
	Description() string
}
