package users

import (
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// JWT Auth https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/

func SignUp(c *gin.Context) {
	// https://www.sohamkamani.com/blog/2018/02/25/golang-password-authentication-and-storage/
	var newUser types.User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unique, err := isNewUserUnique(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email is already taken"})
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into users (username, email, password) values(?,?,?);")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newUser.Username, newUser.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

}

func SignIn(c *gin.Context) {
	var userAuth types.User
	err := c.BindJSON(&userAuth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if userAuth.Password == "" || userAuth.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Password is empty"})
		return
	}

	authenticated, err := auth.IsAuthed(userAuth.Username, userAuth.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	if authenticated {
		token, err := auth.GetToken(userAuth.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create JWT token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login incorrect"})
		return
	}
}

func isNewUserUnique(newUser types.User) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("select COUNT(*) FROM users where username = ? OR email = ?", newUser.Username, newUser.Email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
