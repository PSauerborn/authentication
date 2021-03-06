package main

import (
    "fmt"
    "time"
    "errors"
    "strings"
    "net"
    "regexp"
    "golang.org/x/crypto/bcrypt"
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    log "github.com/sirupsen/logrus"
)

var (
    emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// function used to hash and salt user passwords
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
    user, err := GetUser(PostgresMiddleware{}.Persistence(ctx), username)
    if err != nil {
        return false
    }
    // compare given password with hashed password in database
    log.Info(fmt.Sprintf("checking credentials for user %v+", user))
    return comparePasswords(password, user.Password)
}

// helper function used to determine is a username is already taken
func isUsernameTaken(ctx *gin.Context, username string) bool {
    log.Debug(fmt.Sprintf("checking username %s", username))
    user, err := GetUsername(PostgresMiddleware{}.Persistence(ctx), username)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user from database: %v", err))
        return true
    }
    log.Debug(fmt.Sprintf("username query returned %s", user))
    return user != ""
}

// helper function used to determine is a user email is already taken
func isEmailTaken(ctx *gin.Context, email string) bool {
    log.Debug(fmt.Sprintf("checking email %s", email))
    userEmail, err := GetUserEmail(PostgresMiddleware{}.Persistence(ctx), email)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user email from database: %v", err))
        return true
    }
    log.Debug(fmt.Sprintf("email query returned %s", userEmail))
    return userEmail != ""
}

// function used to generate JWToken with UID and expiry date
func GenerateJWToken(uid string) (string, error) {
    // evaluate expiry time
    expiry := time.Now().UTC()
    expiry = expiry.Add(time.Duration(TokenExpiryMinutes) * time.Minute)

    // generate token and sign with secret key
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "uid": uid,
        "exp": expiry.Unix(),
    })
    return token.SignedString([]byte(JWTSecret))
}

// function used to parse JWT token
func ParseJWToken(tokenString string) (*JWTClaims, error) {
    log.Info(fmt.Sprintf("parsing JWToken %s", tokenString))
    // parse token using JWT secret
    token, _ := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(JWTSecret), nil
    })
    // parse token into custom claims object
    if customClaims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return customClaims, nil
    } else {
        log.Error("unable to parse JWT claims")
        return nil, errors.New("invalid JWToken")
    }
}

// function used to check if email address
func isValidEmail(email string) bool {
    if len(email) < 3 && len(email) > 254 {
        return false
    }
    // check that regex matches
    if !emailRegex.MatchString(email) {
        return false
    }
    parts := strings.Split(email, "@")
    // execute MX record lookup to validate
    mx, err := net.LookupMX(parts[1])
    if err != nil || len(mx) == 0 {
        return false
    }
    return true
}