package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/routes"
)

func main() {
	r := gin.Default()
	r.GET("/", routes.Home)
	r.GET("/a", routes.Home)
	r.GET("/create-url/:url", routes.CreateURL)
	r.GET("/:hash", routes.Redirect)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
