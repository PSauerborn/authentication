package main

import (
    "fmt"
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v4"
    log "github.com/sirupsen/logrus"
)

var (
    PostgresConnection = OverrideStringVariable("POSTGRES_CONNECTION", "postgres://postgres:monty-dev@localhost:5432/authentication")
)


// function used to retrieve user from database
func GetUser(db *pgx.Conn, username string) (User, error) {
    var (user, password string; uid uuid.UUID)
    // get results from database and scan into variables
    results := db.QueryRow(context.Background(), "SELECT uid,username,password FROM users WHERE username=$1", username)
    err := results.Scan(&uid, &user, &password)
    if err != nil {
        log.Error(err)
        return User{}, err
    }
    return User{ Uid: uid, Username: user, Password: password }, nil
}

// function used to retrieve user details from database
func GetUserDetails(db *pgx.Conn, username string) (UserDetails, error) {
    var (uid uuid.UUID; email string; created time.Time; admin bool)
    // get results from database and scan into variables
    results := db.QueryRow(context.Background(), "SELECT uid,email,created,admin FROM user_details WHERE username=$1", username)
    err := results.Scan(&uid, &email, &created, &admin)
    if err != nil {
        log.Error(err)
        return UserDetails{}, err
    }
    return UserDetails{ Uid: uid, Email: email, Created: created, Admin: admin }, nil
}

// function used to retrieve full user from database
func GetFullUser(db *pgx.Conn, username string) (FullUserDetails, error) {
    var (uid uuid.UUID; email string; created time.Time; admin bool)
    // get results from database and scan into variables
    results := db.QueryRow(context.Background(), "SELECT users.uid,email,created,admin FROM user_details details INNER JOIN users users ON (users.uid = details.uid) WHERE username=$1", username)
    err := results.Scan(&uid, &email, &created, &admin)
    if err != nil {
        log.Error(err)
        return FullUserDetails{}, err
    }
    return FullUserDetails{ Uid: uid, Username: username, Email: email, Created: created, Admin: admin }, nil
}

// function used to create a new user in the postgres server
func CreateUser(db *pgx.Conn, user NewUserRequest) error {
    userId := uuid.New()
    log.Info(fmt.Sprintf("insert new user with ID %s", userId))

    // insert values into users table
    _, err := db.Exec(context.Background(), "INSERT INTO users(uid,username,password) VALUES($1,$2,$3)", userId, user.Username, hashAndSalt(user.Password))
    if err != nil {
        log.Error(fmt.Errorf("unable to insert values into users table: %v", err))
        return err
    }
    // insert user into user details table
    _, err = db.Exec(context.Background(), "INSERT INTO user_details(uid,email,admin) VALUES($1,$2,$3)", userId, user.Email, user.Admin)
    if err != nil {
        log.Error(fmt.Errorf("unable to insert values into users_details table: %v", err))
        return err
    }
    return nil
}

// function used to retrieve user email from database
func GetUserEmail(db *pgx.Conn, email string) (string, error) {
    var userEmail string
    results := db.QueryRow(context.Background(), "SELECT email FROM user_details WHERE email=$1", email)
    err := results.Scan(&userEmail)
    if err != nil {
        if err == pgx.ErrNoRows {
            return "", nil
        } else {
            return "", err
        }
    }
    return userEmail, nil
}

// function used to retrieve user username from database
func GetUsername(db *pgx.Conn, username string) (string, error) {
    var user string
    results := db.QueryRow(context.Background(), "SELECT username FROM users WHERE username=$1", username)
    err := results.Scan(&user)
    if err != nil {
        if err == pgx.ErrNoRows {
            return "", nil
        } else {
            return "", err
        }
    }
    return user, nil
}