package main

import (
    "github.com/gin-gonic/gin"
)

// define interface used to store a collection of standard HTTP responses
type StandardHTTPResponse interface{
    Success(ctx *gin.Context)
    InvalidRequestBody(ctx *gin.Context)
    InvalidJSON(ctx *gin.Context)
    InvalidRequest(ctx *gin.Context)
    NotFound(ctx *gin.Context)
    Unauthorized(ctx *gin.Context)
    Forbidden(ctx *gin.Context)
    InternalServerError(ctx *gin.Context)
}

// define set of standard HTTP Responses in JSON format
type StandardJSONResponse struct{}

func(response StandardJSONResponse) Success(ctx *gin.Context) {
    ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "message": "success" })
}

func(response StandardJSONResponse) InvalidRequestBody(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(400, gin.H{ "http_code": 400, "success": false, "message": "invalid request body" })
}

func(response StandardJSONResponse) InvalidJSON(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(400, gin.H{ "http_code": 400, "success": false, "message": "request body must be JSON serializable" })
}

func(response StandardJSONResponse) InvalidRequest(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(400, gin.H{ "http_code": 400, "success": false, "message": "invalid request" })
}

func(response StandardJSONResponse) Unauthorized(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(401, gin.H{ "http_code": 401, "success": false, "message": "unauthorized" })
}

func(response StandardJSONResponse) Forbidden(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(403, gin.H{ "http_code": 403, "success": false, "message": "access forbidden" })
}

func(response StandardJSONResponse) InternalServerError(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(500, gin.H{ "http_code": 500, "success": false, "message": "internal server error" })
}