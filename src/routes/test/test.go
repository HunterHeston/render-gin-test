package test

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/clients"
)

func TestHandler(ctx *gin.Context) {

	client, err := clients.FirestoreClient(ctx)
	if err != nil {
		fmt.Println("error getting firestore client: ", err)
		return
	}

	ref, res, err := client.Collection("test-collection").Add(ctx, map[string]interface{}{
		"TestData": 1234,
	})

	fmt.Println("ref: ", ref)
	fmt.Println("res: ", res)
	fmt.Println("err: ", err)

	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": "The test ran, check the logs...",
	})
}
