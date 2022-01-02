package redirect

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/valuestore"
)

var vs valuestore.ValueStore

// Save a ref to the value store initialized elsewhere.
func New(valueStore valuestore.ValueStore) func(c *gin.Context) {
	vs = valueStore
	return Redirect
}

func Redirect(c *gin.Context) {
	valueID := c.Param("id")
	url, err := vs.LookUp(c, valueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting url for %q", valueID),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"url": string(url),
	})
}
