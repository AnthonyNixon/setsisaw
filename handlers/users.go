package handlers

import (
	"database/sql"
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetCurrentUser(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get user info.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to get sets for themselves.
	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	user_id, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not convert user_id to int, " + err.Error()})
		return
	}

	result := db.QueryRow(database.GET_SPECIFIC_USER, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user info from database, " + err.Error()})
		return
	}

	var user types.User

	err = result.Scan(&user.Id, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetSpecificUser(c *gin.Context) {
	id := c.Param("id")

	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if id != claims.Id {
		// The user is trying to access another person, check if they're an editor or above.
		if !auth.IsEntitled(claims, "EDITOR") {
			c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get user info for user ID %s", claims.Username, id)})
			return
		}
	} else {
		if !auth.IsEntitled(claims, "USER") {
			c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get user info.", claims.Username)})
			return
		}
	}

	// if we made it here, we're good to go.
	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	result := db.QueryRow(database.GET_SPECIFIC_USER, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user info from database, " + err.Error()})
		return
	}

	var user types.User

	err = result.Scan(&user.Id, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "EDITOR") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("User %s is not entitled to get all users.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to get all users.

	user := types.User{}
	users := make([]types.User, 0)

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query(database.GET_ALL_USERS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{"users": users, "count": len(users)})
}
