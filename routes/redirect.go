package routes

import "github.com/gin-gonic/gin"

func Redirect(c *gin.Context) {
	c.JSON(200, gin.H{
		"redirect": "redirect",
	})
}
