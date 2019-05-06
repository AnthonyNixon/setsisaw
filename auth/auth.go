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
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var JWT_SIGNING_KEY []byte
const TOKEN_VALID_TIME = 5 * time.Minute

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

func NewToken(username string, role string) (string, types.Error) {
	var jwtKey = JWT_SIGNING_KEY

	expirationTime := time.Now().Add(TOKEN_VALID_TIME)
	claims := &types.Claims{
		Username: username,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", customerrors.New(http.StatusInternalServerError, err.Error())
	}

	return tokenString, nil
}

func RefreshToken(c *gin.Context) (string, types.Error) {
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

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > time.Minute {
		return "", customerrors.New(http.StatusBadRequest, "too early to refresh token, token is valid for more than 1 minute.")
	}

	expirationTime := time.Now().Add(TOKEN_VALID_TIME)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JWT_SIGNING_KEY)
	if err != nil {
		return "", customerrors.New(http.StatusInternalServerError, err.Error())
	}

	return tokenString, nil
}

func GetUserInfo(c *gin.Context) (types.Claims, types.Error) {
	authHeader := c.GetHeader("Authorization")
	claims := &types.Claims{}

	tokenString, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		return *claims, customerrors.New(http.StatusInternalServerError, err.Error())
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SIGNING_KEY, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return *claims, customerrors.New(http.StatusUnauthorized, "signature invalid, " + err.Error())
		}

		if !tkn.Valid {
			return *claims, customerrors.New(http.StatusUnauthorized, err.Error())
		}

		return *claims, customerrors.New(http.StatusBadRequest, "Invalid JWT token, " + err.Error())

	}

	return *claims, nil
}

func IsEntitled(claims types.Claims, requirement string) bool {
	switch strings.ToUpper(claims.Role) {
	case "USER":
		switch requirement {
		case "USER":
			return true
		default:
			return false
		}
	case "EDITOR":
		switch requirement {
		case "USER", "EDITOR":
			return true
		default:
			return false
		}
	case "ADMIN":
		return true
	default:
		return false
	}
}