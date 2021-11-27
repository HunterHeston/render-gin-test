package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Redirect(c *gin.Context) {

	ds := dataStore{
		database,
	}

	hash := c.Param("hash")
	url, err := ds.lookupURL(hash)
	if err != nil {
		fmt.Printf("could not find url for hash %q\n", hash)
	}

	fmt.Printf("hash is %q\n", hash)

	c.JSON(200, gin.H{
		"url": url,
	})
}
