package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	flagList = flag.StringP(
		"list",
		"l",
		"",
		"list existing credentials based on description",
	)
	flagNew = flag.StringP(
		"new",
		"n",
		"",
		"generate new credentials with human readable description",
	)
	flagRemove = flag.StringP(
		"remove",
		"r",
		"",
		"Remove credentials by id",
	)
	flagMongo = flag.StringP(
		"mongo-uri",
		"m",
		"mongodb://localhost:27017",
		"MongoDB instance to connect to (required)",
	)
)

func main() {
	flag.Parse()
	client := connectDb()
	c := client.Database("credentials").Collection("clients")
	switch {
	case *flagNew != "":
		create(c)
	case *flagRemove != "":
		remove(c)
	default:
		list(c)
	}
}

type Creds struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Description  string `json:"description"`
}

func create(c *mongo.Collection) {

	p, err := randomHex(32)
	if err != nil {
		log.Error(err)
		return
	}
	cred := Creds{
		ClientId:     uuid.New().String(),
		ClientSecret: p,
		Description:  *flagNew,
	}
	b, err := json.Marshal(cred)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	en := encrypt(p)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = c.InsertOne(
		ctx,
		bson.D{{"clientId", cred.ClientId},
			{"clientSecret", string(en)},
			{"description", cred.Description}},
	)
	if err != nil {
		log.Error(err)
	}
	return
}

func encrypt(s string) []byte {
	b := []byte(s)
	en, err := bcrypt.GenerateFromPassword(b, 12)
	if err != nil {
		fmt.Println(err)
		n := make([]byte, 0)
		return n
	}
	return en
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func list(c *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{}
	if *flagList != "" {
		filter = bson.D{{"description", *flagList}}
	}
	cur, err := c.Find(
		ctx,
		filter,
		options.Find().SetProjection(bson.D{{"clientSecret", 0}, {"_id", 0}}),
	)
	if err != nil {
		log.Error(err)
		return
	}
	var creds []Creds
	for cur.Next(ctx) {
		var cred Creds
		err := cur.Decode(&cred)
		if err != nil {
			log.Error(err)
		}
		creds = append(creds, cred)
	}
	b, err := json.Marshal(creds)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

// open connection to mongo db
func connectDb() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongo))
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Error(err)
	}
	return client
}

func remove(c *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.DeleteOne(
		ctx,
		bson.D{{"clientId", *flagRemove}},
	)
	if err != nil {
		log.Error(err)
	}
	log.Info(res)
}
