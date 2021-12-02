package routes

import (
	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Nothing to see here.",
	})
}
