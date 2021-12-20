package firestore

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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
		return nil, fmt.Errorf("error getting doc %q: %v", id, err)
	}

	url, ok := doc.Data()["url"]
	if !ok {
		return nil, fmt.Errorf("error document data does not contain url")
	}

	strURL, ok := url.(string)
	if !ok {
		return nil, fmt.Errorf("error url data is not of type string")
	}

	return []byte(strURL), nil
}

///////////////////////////////////////////
func (f Firestore) Save(value []byte) (string, error) {
	return "", nil
}
