package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// main route handler
func CreateURL(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"shortURL": "",
	})
}

type dataStore struct {
	store map[string]string
}
