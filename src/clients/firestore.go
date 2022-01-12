package clients

import (
	"context"
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	KEY_PATH_ARG = "FIREBASE_SERVICE_KEY_PATH"
)

var (
	serviceKeyPath string
)

func init() {
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

func FirestoreClient(ctx context.Context) (*firestore.Client, error) {
	// Use a service account
	sa := option.WithCredentialsFile(serviceKeyPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
