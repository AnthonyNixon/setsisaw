package main

import (
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/handlers"
	"github.com/AnthonyNixon/setsisaw/users"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var PORT = ""

func init() {
	database.Initialize()
	auth.Initialize()
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
}

func main() {
	r := gin.Default()

	r.POST("/signup", users.SignUp)

	r.POST("/signin", users.SignIn)

	r.GET("/authcheck", handlers.AuthCheck)

	r.GET("/refresh", handlers.RefreshToken)

	// Artists
	r.POST("/artists", handlers.NewArtist)

	// Locations
	r.POST("/locations", handlers.NewLocation)

	// Sets
	r.POST("/sets", handlers.NewSet)


	log.Printf("Running SetsISaw API on :%s...", PORT)
	err := r.Run(fmt.Sprintf(":%s", PORT)) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal(err.Error())
	}
}
