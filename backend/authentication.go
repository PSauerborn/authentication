package main

import (
    "fmt"
    "github.com/jackc/pgx/v4"
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    jaeger_negroni "github.com/PSauerborn/jaeger-negroni"
)

var (
    StandardHTTP = StandardJSONResponse{}
)

// function used to start new authentication service
func main() {
    // read environment variables from config into local variables
    ConfigureService()
    router := gin.New()
    router.Use(CORSMiddleware())

    // create new config object
    jaegerConfig := jaeger_negroni.Config("jaeger-agent", "authentication-api", 6831)
    tracer := jaeger_negroni.NewTracer(jaegerConfig)
    defer tracer.Close()

    router.Use(jaeger_negroni.JaegerNegroni(jaegerConfig))

    // configure GET routes used for server
    router.GET("/health", HealthCheckHandler)
    // configure POST routes used for server
    router.POST("/token", PostgresMiddleware{}.Middleware(), GetTokenHandler)
    router.POST("/signup", PostgresMiddleware{}.Middleware(), CreateUserHandler)
    router.POST("/user", PostgresMiddleware{}.Middleware(), GetUserDetailsHandler)

    log.Info(fmt.Sprintf("starting authentication service at %s:%d", ListenAddress, ListenPort))
    router.Run(fmt.Sprintf("%s:%d", ListenAddress, ListenPort))
}

// function used as basic health check
func HealthCheckHandler(context *gin.Context) {
    StandardHTTP.Success(context)
}

// function used to retrieve JWToken for user. Note that both
// username and password must be present in the request body
// in order for the JWT token to be generated successfully
func GetTokenHandler(ctx *gin.Context) {
    var request TokenRequest
    err := ctx.ShouldBind(&request)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse request body: %v", err))
        StandardHTTP.InvalidRequestBody(ctx)
        return
    }
    log.Info(fmt.Sprintf("received request for token from user %s", request.Uid))
    // if user is not authorized, return 401 response
    if (!isAuthenticatedUser(ctx, request.Uid, request.Password)) {
        log.Warn(fmt.Sprintf("received invalid login request for user %s", request.Uid))
        StandardHTTP.Unauthorized(ctx)
        return
    }
    // generate JWToken if user is successfully authenticated and return in response
    log.Info(fmt.Sprintf("successfully authenticated user %s. Generating JWT", request.Uid))
    token, err := GenerateJWToken(request.Uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to create JWT: %s", err))
        StandardHTTP.InternalServerError(ctx)
    } else {
        log.Debug(fmt.Sprintf("generated token %s", token))
        ctx.JSON(200, gin.H{ "http_code": 200, "token": token })
    }
}

// function used to create new user in database. username, email and password
// must be present in the request body. Note that all passwords are hashed
// and salted before they are stored in the database
func CreateUserHandler(ctx *gin.Context) {
    var request NewUserRequest
    err := ctx.ShouldBind(&request)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse request body: %v", err))
        StandardHTTP.InvalidRequestBody(ctx)
        return
    }
    // check if username is already taken
    if (isUsernameTaken(ctx, request.Username)) {
        payload := gin.H{"http_code": 400, "success": false, "message": "username already in use"}
        ctx.AbortWithStatusJSON(400, payload)
        return
    }
    // check that email address is valid
    if (!isValidEmail(request.Email)) {
        log.Error(fmt.Sprintf("received request with invalid email address %s", request.Email))
        payload := gin.H{"http_code": 400, "success": false, "message": "invalid email address"}
        ctx.AbortWithStatusJSON(400, payload)
        return
    }
    // check if email is already taken
    if (isEmailTaken(ctx, request.Email)) {
        payload := gin.H{"http_code": 400, "success": false, "message": "email already in use"}
        ctx.AbortWithStatusJSON(400, payload)
        return
    }
    // insert new values into databawse
    err = CreateUser(ctx.MustGet("persistence").(*pgx.Conn), request)
    if err != nil {
        StandardHTTP.InternalServerError(ctx)
    } else {
        log.Debug("successfully created new user")
        ctx.JSON(200, gin.H{ "http_code": 200, "message": fmt.Sprintf("successfully created new user '%s'", request.Username) })
    }
}

// function used to retrieve details for user
func GetUserDetailsHandler(ctx *gin.Context) {
    var request IntrospectionRequest
    err := ctx.ShouldBind(&request)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse request body: %v", err))
        StandardHTTP.InvalidRequestBody(ctx)
        return
    }
    log.Debug(fmt.Sprintf("retrieved body %+v", request))
    claims, err := ParseJWToken(request.Token)
    if err != nil {
        StandardHTTP.Unauthorized(ctx)
    } else {
        user, err := GetFullUser(PostgresMiddleware{}.Persistence(ctx), claims.Uid)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve user details: %v", err))
            StandardHTTP.InternalServerError(ctx)
            return
        }
        ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": user })
    }
}