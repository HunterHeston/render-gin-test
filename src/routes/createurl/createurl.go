package createurl

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/validation"
	"github.com/hunterheston/gin-server/src/valuestore"
)

var vs valuestore.ValueStore

func New(valueStore valuestore.ValueStore) func(c *gin.Context) {
	vs = valueStore
	return CreateURL
}

// main route handler
func CreateURL(c *gin.Context) {

	fmt.Println("HSH create: ", c.ClientIP())
	fip := c.Request.Header["X-Forwarded-For"]
	fmt.Println("HSH create X-Forwarded-For: ", fip)
	fmt.Println("HSH create Header: ", c.Request.Header)

	// route expects a single param "url"
	rawUrlInput := c.Query("url")

	rawUrlInput = appendHTTP(rawUrlInput)

	// must be a valid url.
	if !validation.ValidateURL(rawUrlInput) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Input must be a valid url, %q is not valid.", rawUrlInput),
		})
		return
	}

	// Store the url and get it's id.
	valueID, err := vs.Save(c, []byte(rawUrlInput))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Error saving your input: %v", err),
		})
	}

	frontEnd := os.Getenv("FRONTEND_HOST")

	// Send the ID of the stored value back to the client.
	c.JSON(http.StatusOK, gin.H{
		"id":  valueID,
		"url": fmt.Sprintf("%s/%s", frontEnd, valueID),
	})
}

func appendHTTP(url string) string {
	if !strings.Contains(url, "http") {
		return "http://" + url
	}
	return url
}
