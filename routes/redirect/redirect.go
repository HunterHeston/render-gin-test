package redirect

import (
	"fmt"

	"github.com/gin-gonic/gin"
	valuestore "github.com/hunterheston/gin-server/helpers/valuestore"
)

var vs valuestore.ValueStore

// Save a ref to the value store initialized elsewhere.
func New(valueStore valuestore.ValueStore) func(c *gin.Context) {
	vs = valueStore
	return Redirect
}

func Redirect(c *gin.Context) {
	key := c.Param("hash")
	url, err := vs.LookUp(key)
	if err != nil {
		fmt.Printf("Error looking up value for key=%v: %v", key, err)
	}
	c.JSON(200, gin.H{
		"url": url,
	})
}
