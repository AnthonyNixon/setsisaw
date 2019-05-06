package users

import (
	"database/sql"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/customerrors"
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
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind."})
		return
	}

	if newUser.Email == "" || newUser.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must include email and username"})
		return
	}

	unique, err := isNewUserUnique(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error3": "Username or email is already taken"})
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error4": err.Error()})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error5": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into users (username, email, password, first_name, last_name) values(?,?,?,?,?);")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newUser.Username, newUser.Email, hashedPassword, newUser.FirstName, newUser.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"username": newUser.Username, "email": newUser.Email})
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
		role, err := getUserRole(userAuth.Username)
		if err != nil {
			c.JSON(err.StatusCode(), gin.H{"error": err.Description()})
		}

		token, err := auth.NewToken(userAuth.Username, role)
		if err != nil {
			c.JSON(err.StatusCode(), gin.H{"error": err.Description()})
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

func getUserRole(username string) (string, types.Error) {
	db, err := database.GetConnection()
	if err != nil {
		return "", customerrors.New(http.StatusInternalServerError, "could not connect to database")
	}
	defer db.Close()

	var role string
	result := db.QueryRow("select role FROM users where username = ?", username)
	if err != nil {
		return "", customerrors.New(http.StatusInternalServerError, "could not query database")
	}

	err = result.Scan(&role)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			return "", customerrors.New(http.StatusInternalServerError, "could not find user role in database")
		}

		return "", customerrors.New(http.StatusInternalServerError, "could not query database")
	}

	return role, nil
}
