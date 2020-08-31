package main

import (
	"fmt"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PostgresMiddleware struct{}

// define middleware used to instantiate postgres connection
// when request reaches server. The Postgres connection is then
// set as an attribute of the Gin/Gonic context
func (middleware PostgresMiddleware) Middleware() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		// create database connection. return internal server error if connection fails
		db, err := pgx.Connect(context.Background(), PostgresConnection)
		if err != nil {
			log.Error(fmt.Errorf("unable to connect to postgres server: %v", err))
			StandardHTTP.InternalServerError(ctx)
		} else {
			// defer closing of connection and set persistence as attribute of context
			defer db.Close(context.Background())
			log.Debug("successfully connected to postgres server. setting as attribute of context")
			ctx.Set("persistence", db)
			ctx.Next()
		}
	}
}

// function used to retrieve persistence from context
func (middleware PostgresMiddleware) Persistence(ctx *gin.Context) *pgx.Conn {
	return ctx.MustGet("persistence").(*pgx.Conn)
}

// function used to add CORS headers to routes
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET,PUT,PATCH,DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}