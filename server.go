package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/routes/createurl"
	"github.com/hunterheston/gin-server/src/routes/redirect"
	"github.com/hunterheston/gin-server/src/routes/root"
	valuestore "github.com/hunterheston/gin-server/src/valuestore"
	"github.com/hunterheston/gin-server/src/valuestore/inmemory"
	"github.com/joho/godotenv"
)

var (
	database       valuestore.ValueStore
	allowedOrigins = []string{"http://localhost:3000"}
)

func init() {
	// setup the data store used throughout the server.
	database = inmemory.NewInMemory()

	// load and save env vars
	godotenv.Load()
	frontendHost := os.Getenv("FRONTEND_HOST")
	allowedOrigins = append(allowedOrigins, frontendHost)
	fmt.Println("HSH ", allowedOrigins)
}

func main() {
	r := gin.Default()

	// not using any proxies currently and defaulting to trusting all is insecure
	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	// define hosts that can make requests to this server.
	r.Use(cors.New(
		cors.Config{
			AllowOrigins: allowedOrigins,
		},
	))

	r.GET("/", root.Root)
	r.GET("/new", createurl.New(database))
	r.GET("/:id", redirect.New(database))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
