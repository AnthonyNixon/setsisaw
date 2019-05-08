package types

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type Artist struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	DefaultGenre string `json:"default_genre"`
}

type SetMetadata struct {
	Rating int    `json:"rating"`
	Genre  string `json:"genre"`
	Length int    `json:"length"`
	Notes  string `json:"notes"`
}

type Set struct {
	Id           int         `json:"id"`
	UserId       int         `json:"user_id"`
	ArtistId     int         `json:"artist_id"`
	ArtistName   string      `json:"artist_name"`
	LocationId   int         `json:"location_id"`
	LocationName string      `json:"location_name"`
	Date         string      `json:"date"`
	Metadata     SetMetadata `json:"metadata"`
}

type Location struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	IsFestival  bool   `json:"is_festival"`
	Year        int    `json:"year"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Id       string `json:"id"`
	jwt.StandardClaims
}

type Error interface {
	StatusCode() int
	Description() string
}
