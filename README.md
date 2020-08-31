### Identity Provider


Simple, RESTful Identity Provider that provides a route to retrieved JWTokens used by applications sitting on
my personal VM. The application itself is written in `Go` and uses a `PostgreSQL` database to store
user credentials. All responses are returned in JSON format; the API is designed to operate as a microservice
and can easily deployed with the provided `Dockerfile`. See `swagger.yaml` documentation provided
to see available routes.

Note that all user credentials are hashed and salted before they are stored in the `Postgres` server