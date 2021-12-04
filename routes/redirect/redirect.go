package redirect

import (
	"github.com/gin-gonic/gin"
	valuestore "github.com/hunterheston/gin-server/helpers/valuestore"
)

var vs valuestore.ValueStore

// Save a ref to the value store initialized elsewhere.
func NewRedirect(valueStore valuestore.ValueStore) func(c *gin.Context) {
	vs = valueStore
	return Redirect
}

func Redirect(c *gin.Context) {
	key := c.Param("hash")
	url := vs.LookUp(key)
	c.JSON(200, gin.H{
		"url": url,
	})
}
