package root

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "There's nothing here.",
	})
}
