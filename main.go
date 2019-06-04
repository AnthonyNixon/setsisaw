package main

import (
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/handlers"
	"github.com/AnthonyNixon/setsisaw/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
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
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/signup", users.SignUp)
	r.POST("/signin", users.SignIn)
	r.GET("/refresh", handlers.RefreshToken)

	r.GET("/authcheck", handlers.AuthCheck)

	// Users
	r.GET("/users", handlers.GetAllUsers)
	r.GET("/user/current", handlers.GetCurrentUser)
	r.GET("/users/:id", handlers.GetSpecificUser)
	r.PUT("/users", handlers.UpdateUser)

	// Artists
	r.POST("/artists", handlers.NewArtist)
	r.GET("/artists", handlers.GetAllArtists)
	r.GET("/artists/:id", handlers.GetArtist)

	// Locations
	r.POST("/locations", handlers.NewLocation)
	r.GET("/locations", handlers.GetAllLocations)
	r.GET("/locations/:id", handlers.GetLocation)
	r.PUT("/locations/:id", handlers.UpdateLocation)

	// Sets
	r.POST("/sets", handlers.NewSet)
	r.GET("/sets", handlers.GetSetsForCurrentUser)
	r.GET("/sets/all", handlers.GetAllSets)

	log.Printf("Running SetsISaw API on :%s...", PORT)

	err := r.Run(fmt.Sprintf(":%s", PORT)) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal(err.Error())
	}
}
