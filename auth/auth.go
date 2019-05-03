package auth

import (
	"database/sql"
	"github.com/AnthonyNixon/setsisaw/customerrors"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
	"github.com/AnthonyNixon/setsisaw/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var JWT_SIGNING_KEY []byte

func Initialize() {
	log.Print("Initializing Authentication")
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("No Signing Key Present.")
	}

	JWT_SIGNING_KEY = []byte(signingKey)
	log.Print("done")
}

func IsAuthed(username string, password string) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var storedPassword string
	result := db.QueryRow("select password FROM users where username = ?", username)
	if err != nil {
		return false, err
	}

	err = result.Scan(&storedPassword)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, nil
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func GetToken(username string) (string, error) {
	var jwtKey = JWT_SIGNING_KEY

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &types.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GetUsernameFromAuthHeader(c *gin.Context) (string, types.Error) {
	authHeader := c.GetHeader("Authorization")
	tokenString, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		return "", customerrors.New(http.StatusInternalServerError, err.Error())
	}

	claims := &types.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SIGNING_KEY, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", customerrors.New(http.StatusUnauthorized, "signature invalid, " + err.Error())
		}

		if !tkn.Valid {
			return "", customerrors.New(http.StatusUnauthorized, err.Error())
		}

		return "", customerrors.New(http.StatusBadRequest, "Invalid JWT token, " + err.Error())

	}

	return claims.Username, nil
}
