package main

import (
	"time"
	"github.com/google/uuid"
)


// define struct used to hold JWT Claims
type JWTClaims struct {
	Uid   string
	Admin bool
	Exp   time.Time
}

// ###################################################
// # Define structs used to decode JSON request bodies
// ###################################################

// define struct used to define request body for new token
type TokenRequest struct {
	Uid      string `json:"uid" binding:"required"`
	Password string	`json:"password" binding:"required"`
}

// define struct used to define request body for new token
type NewUserRequest struct {
	Username string `json:"uid" binding:"required"`
	Password string	`json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Admin    *bool  `json:"admin" binding:"required"`
}

// #######################################################
// # Define structs used for persistence responses from db
// #######################################################

// define struct used to encapsulate user details
type User struct {
	Uid      uuid.UUID `json:"uid"`
	Username string	   `json:"username"`
	Password string    `json:"password"`
}

// define struct used to encapsulate user details
type UserDetails struct {
	Uid     uuid.UUID `json:"uid"`
	Email   string	  `json:"email"`
	Created time.Time `json:"created"`
	Admin   bool      `json:"admin"`
}

type FullUserDetails struct {
	Uid      uuid.UUID `json:"uid"`
	Username string	   `json:"username"`
	Email    string	   `json:"email"`
	Created  time.Time `json:"created"`
	Admin    bool      `json:"admin"`
}