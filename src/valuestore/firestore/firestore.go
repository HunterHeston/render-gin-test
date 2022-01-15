package firestore

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/hunterheston/gin-server/src/clients"
	"github.com/hunterheston/gin-server/src/stringgeneration"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
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

		counter := Counter{
			numShards: 10,
		}

		if err := counter.initCounter(ctx, f.c.Doc(docPath)); err != nil {
			return "", err
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

	counter := Counter{
		numShards: 10,
	}

	if _, err := counter.incrementCounter(ctx, docRef); err != nil {
		return err
	}

	count, _ := counter.getCount(ctx, docRef)
	fmt.Println("HSH count: ", count)

	return nil
}

////////////////////////////////
// distributed counter
////////////////////////////////
// Counter is a collection of documents (shards)
// to realize counter with high frequency.
type Counter struct {
	numShards int
}

// Shard is a single counter, which is used in a group
// of other shards within Counter.
type Shard struct {
	Count int
}

// initCounter creates a given number of shards as
// subcollection of specified document.
func (c *Counter) initCounter(ctx context.Context, docRef *firestore.DocumentRef) error {
	colRef := docRef.Collection("shards")

	// Initialize each shard with count=0
	for num := 0; num < c.numShards; num++ {
		shard := Shard{0}

		if _, err := colRef.Doc(strconv.Itoa(num)).Set(ctx, shard); err != nil {
			return fmt.Errorf("Set: %v", err)
		}
	}
	return nil
}

// incrementCounter increments a randomly picked shard.
func (c *Counter) incrementCounter(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	docID := strconv.Itoa(rand.Intn(c.numShards))

	shardRef := docRef.Collection("shards").Doc(docID)
	return shardRef.Update(ctx, []firestore.Update{
		{Path: "Count", Value: firestore.Increment(1)},
	})
}

// getCount returns a total count across all shards.
func (c *Counter) getCount(ctx context.Context, docRef *firestore.DocumentRef) (int64, error) {
	var total int64
	shards := docRef.Collection("shards").Documents(ctx)
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("Next: %v", err)
		}

		vTotal := doc.Data()["Count"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, fmt.Errorf("firestore: invalid dataType %T, want int64", vTotal)
		}
		total += shardCount
	}
	return total, nil
}
