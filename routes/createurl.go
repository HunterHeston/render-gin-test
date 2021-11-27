package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var database = make(map[string]string)

// main route handler
func CreateURL(c *gin.Context) {
	// ds := dataStore{
	// 	store: database,
	// }

	// c.Request.Body
	// ds.createurl()
}

type dataStore struct {
	store map[string]string
}

func (d dataStore) createurl(url string) (string, error) {
	if d.store == nil {
		return "", fmt.Errorf("data store is nil")
	}

	if !validateURL(url) {
		return "", fmt.Errorf("url %q is invalid", url)
	}
	return "", nil
}

func validateURL(url string) bool {

	return true
}
