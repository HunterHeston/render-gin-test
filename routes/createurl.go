package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var database = make(map[string]string)

// main route handler
func CreateURL(c *gin.Context) {
	ds := dataStore{
		store: database,
	}

	fmt.Printf("Data store: %v\n", ds)

	url := c.Param("url")
	shortURL, err := ds.createurl(url)
	if err != nil {
		fmt.Println("Error getting short url: ", err)
	}

	fmt.Printf("Orignial url: %q short url: %q\n", url, shortURL)

	c.JSON(http.StatusOK, gin.H{
		"shortURL": shortURL,
	})
}

type dataStore struct {
	store map[string]string
}

func (ds dataStore) createurl(url string) (string, error) {
	if ds.store == nil {
		return "", fmt.Errorf("data store is nil")
	}

	if !validateURL(url) {
		return "", fmt.Errorf("url %q is invalid", url)
	}

	hash := rand.Intn(100)
	ds.store[strconv.FormatInt(int64(hash), 10)] = url

	return fmt.Sprintf("localhost:8080/%d", hash), nil
}

func (ds dataStore) lookupURL(hash string) (string, error) {
	return ds.store[hash], nil
}

func validateURL(url string) bool {

	return true
}
