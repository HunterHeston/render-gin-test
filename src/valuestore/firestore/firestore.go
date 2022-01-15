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
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/clients"
	"github.com/hunterheston/gin-server/src/stringgeneration"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	KEY_PATH_ARG                       = "FIREBASE_SERVICE_KEY_PATH"
	FIRESTORE_PATH_URL_COLLECTION_PATH = "FIRESTORE_PATH_URL_COLLECTION_PATH"
)

var (
	serviceKeyPath string
	firestorePath  string
)

///////////////////////////////////////////
func init() {
	godotenv.Load()
	fp := os.Getenv(FIRESTORE_PATH_URL_COLLECTION_PATH)
	if fp == "" {
		fmt.Printf("ERROR: could not get %q from environment.\n", FIRESTORE_PATH_URL_COLLECTION_PATH)
	}
	firestorePath = fp

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
type Firestore struct {
	c *firestore.Client
}

///////////////////////////////////////////
func getClient(ctx context.Context) *firestore.Client {
	// Use a service account
	sa := option.WithCredentialsFile(serviceKeyPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Error setting up new firebase app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting client: %v\n", err)
	}

	return client
}

///////////////////////////////////////////
func New(ctx context.Context) Firestore {
	return Firestore{
		c: getClient(ctx),
	}
}

///////////////////////////////////////////
func (f Firestore) LookUp(ctx context.Context, id string) ([]byte, error) {
	doc, err := f.c.Collection(firestorePath).Doc(id).Get(ctx)
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
	docPath := path.Join(firestorePath, id)
	if err := f.updateLookupMetric(ctx, docPath); err != nil {
		fmt.Printf("error updating metrics for %q: %v", docPath, err)
	}

	return []byte(strURL), nil
}

///////////////////////////////////////////
func (f Firestore) Save(ctx context.Context, value []byte) (string, error) {

	done := false
	attempts := 0
	resultID := ""
	for !done && attempts < 10 {
		attempts++

		id := stringgeneration.RandStringBytesRmndr(6)
		docPath := path.Join(firestorePath, id)
		_, err := f.c.Doc(docPath).Create(ctx, struct{}{})

		// error getting doc. Try again.
		if err != nil {
			fmt.Printf("Error doc already exists doc %q: %v\n", id, err)
			continue
		}

		resultID = id
		_, err = f.c.Doc(docPath).Set(ctx, map[string]interface{}{
			"url":          string(value),
			"date_created": time.Now(),
		})
		// exit
		if err != nil {
			fmt.Println("Error creating new doc: ", err)
		}

		break
	}

	return resultID, nil
}

///////////////////////////////////////////
// Metrics for save and lookup
///////////////////////////////////////////
type lookupMetric struct {
	ip        string
	timestamp time.Time
}

func (f Firestore) updateLookupMetric(ctx context.Context, path string) error {
	fs, err := clients.FirestoreClient(ctx)
	if err != nil {
		return err
	}

	c, ok := ctx.(*gin.Context)
	if !ok {
		return fmt.Errorf("error casting context.Context to gin.Context")
	}

	metric := lookupMetric{}

	metric.ip = c.Request.Header.Get("CF-Connecting-IP")
	metric.timestamp = time.Now()

	docRef := fs.Doc(path)
	docRef.Collection("metrics").Add(ctx, map[string]interface{}{
		"ip":        metric.ip,
		"timestamp": metric.timestamp,
	})
	return nil
}
