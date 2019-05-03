package main

import (
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/handlers"
	"github.com/AnthonyNixon/setsisaw/users"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	database.Initialize()
	auth.Initialize()
}

func main() {
	r := gin.Default()

	r.POST("/signup", users.SignUp)

	r.POST("/signin", users.SignIn)

	r.GET("/authcheck", handlers.AuthCheck)

	log.Print("Running SetsISaw API...")
	r.Run() // listen and serve on 0.0.0.0:8080
}
