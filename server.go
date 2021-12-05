package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/routes/createurl"
	"github.com/hunterheston/gin-server/src/routes/redirect"
	valuestore "github.com/hunterheston/gin-server/src/valuestore"
	"github.com/hunterheston/gin-server/src/valuestore/inmemory"
)

var database valuestore.ValueStore

func init() {
	database = inmemory.NewInMemory()
}

func main() {

	r := gin.Default()
	r.GET("/new", createurl.New(database))
	r.GET("/:hash", redirect.New(database))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
