package firestore

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/hunterheston/gin-server/src/stringgeneration"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const KEY_PATH_ARG = "FIREBASE_SERVICE_KEY_PATH"

var serviceKeyPath string

type Firestore struct {
	c *firestore.Client
}

///////////////////////////////////////////
func init() {
	// load env vars
	godotenv.Load()

	serviceKeyPath = os.Getenv(KEY_PATH_ARG)
	if serviceKeyPath == "" {
		fmt.Printf("Error: firebase service key path not set.\n")
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working dir: %v\n", err)
	}
	serviceKeyPath = path.Join(wd, serviceKeyPath)
}

///////////////////////////////////////////
func getClient() *firestore.Client {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile(serviceKeyPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln("HSH", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("HSH", err)
	}

	return client
}

///////////////////////////////////////////
func New() Firestore {
	return Firestore{
		c: getClient(),
	}
}

///////////////////////////////////////////
func (f Firestore) LookUp(id string) ([]byte, error) {
	ctx := context.Background()
	doc, err := f.c.Collection("urls").Doc(id).Get(ctx)
	if err != nil {
		fmt.Println("error getting doc")
		return nil, fmt.Errorf("error getting doc %q: %v", id, err)
	}

	url, ok := doc.Data()["url"]
	if !ok {
		fmt.Println("error document data does not contain url")
		return nil, fmt.Errorf("error document data does not contain url")
	}

	strURL, ok := url.(string)
	if !ok {
		fmt.Println("error url data is not of type string")
		return nil, fmt.Errorf("error url data is not of type string")
	}

	return []byte(strURL), nil
}

///////////////////////////////////////////
func (f Firestore) Save(value []byte) (string, error) {
	ctx := context.Background()
	done := false
	attempts := 0
	resultID := ""
	for !done && attempts < 10 {
		attempts++

		id := stringgeneration.RandStringBytesRmndr(6)
		doc, err := f.c.Collection("urls").Doc(id).Get(ctx)

		// error getting doc. Try again.
		if status.Code(err) != codes.NotFound {
			fmt.Printf("Error fetching doc %q: %v\n", id, err)
			continue
		}
		// id is already in use. Try again.
		if doc.Exists() {
			continue
		}

		resultID = id
		f.c.Collection("urls").Doc(id).Set(ctx, map[string]interface{}{
			"url":          string(value),
			"date_created": time.Now(),
		})
		// exit
		break
	}

	return resultID, nil
}
