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

func GetClient() *firestore.Client {
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

func GetServiceAccount() {
	//TODO: get the service account key from env vars
	// create a service account instance
	// return the account
}
