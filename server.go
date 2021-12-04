package main

import (
	"github.com/gin-gonic/gin"
	valuestore "github.com/hunterheston/gin-server/helpers/valuestore"
	"github.com/hunterheston/gin-server/helpers/valuestore/inmemory"
	createurl "github.com/hunterheston/gin-server/routes/createurl"
	"github.com/hunterheston/gin-server/routes/redirect"
)

func init() {
	// initialize an in memory data store.
	// make routes initializable.
	// allow passing of value into a route structure.
	// need to pass the data store to a route.
	// So that both CreateURL and Redirect URL can access the same in memory data store.
}

var database valuestore.ValueStore

func init() {
	database = inmemory.NewInMemory()

}

func main() {

	r := gin.Default()
	r.GET("/create-url/:url", createurl.CreateURL)
	r.GET("/:hash", redirect.NewRedirect(database))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
