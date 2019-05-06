package types

import "github.com/dgrijalva/jwt-go"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Id int `json:"id"`
}

type Artist struct {
	Name string `json:"name"`
	Id int `json:"id"`
}

type Claims struct {
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type Error interface {
	StatusCode() int
	Description() string
}
