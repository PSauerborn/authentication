package main

import (
	"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// function used to hash users passwords
func hashAndSalt(password string) string {
	// convert passwords into byte array and hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
	}
	return string(hash)
}

// function used to compare password to hashed password from database
func comparePasswords(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Warn(err)
		return false
	}
	return true
}

// function used to authenticate username
func isAuthenticatedUser(ctx *gin.Context, username, password string) bool {
	// retrieve user details from database and compare to given password
	user, err := GetUser(PostgresMiddleware{}.GetConnection(ctx), username)
	if err != nil {
		return false
	}
	log.Info(fmt.Sprintf("checking credentials for user %v+", user))
	return comparePasswords(password, user.Password)
}

// function used to generate JWToken with UID and expiry date
func GenerateJWToken(uid string) (string, error) {
	// evaluate expiry time
	expiry := time.Now().UTC()
	expiry.Add(time.Duration(TokenExpiryMinutes) * time.Minute)

	// generate token and sign with secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"expiry": expiry,
	})
	return token.SignedString([]byte(JWTSecret))
}

func ParseJWToken(token string) (JWTClaims, error) {
	return JWTClaims{}, nil
}