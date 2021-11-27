package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/routes"
)

func main() {
	r := gin.Default()
	r.GET("/home", routes.Home)
	r.GET("/create-url", routes.CreateURL)
	// r.NoRoute(routes.Redirect)
	r.GET("/:hash", routes.Redirect)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
