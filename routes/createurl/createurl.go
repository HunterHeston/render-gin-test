package createurl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/helpers/valuestore"
)

var vs valuestore.ValueStore

func New(valueStore valuestore.ValueStore) func(c *gin.Context) {
	vs = valueStore
	return CreateURL
}

// main route handler
func CreateURL(c *gin.Context) {

	// route expects a single
	rawUrlInput := c.Query("url")
	if !validUrl(rawUrlInput) {
		// TODO: return an error to the client.
	}

	valueID, err := vs.Save([]byte(rawUrlInput))
	if err != nil {
		// TODO: return an error to the client.
	}

	c.JSON(http.StatusOK, gin.H{
		"id": valueID,
	})
}

// check if string is a valid url.
func validUrl(url string) bool {
	return true
}

type dataStore struct {
	store map[string]string
}
