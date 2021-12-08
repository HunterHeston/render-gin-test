package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/routes/createurl"
	"github.com/hunterheston/gin-server/src/routes/redirect"
	"github.com/hunterheston/gin-server/src/routes/root"
	valuestore "github.com/hunterheston/gin-server/src/valuestore"
	"github.com/hunterheston/gin-server/src/valuestore/inmemory"
)

var database valuestore.ValueStore

func init() {
	database = inmemory.NewInMemory()
}

func main() {

	r := gin.Default()
	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://localhost:3000"},
		},
	))

	r.GET("/", root.Root)
	r.GET("/new", createurl.New(database))
	r.GET("/:id", redirect.New(database))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
