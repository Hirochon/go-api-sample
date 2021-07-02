package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"work/middleware"

	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var addr = ":8080"

type helloJSON struct {
	UserName string `json:"user_name"`
	Content  string `json:"content"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GET hello!\n")
	case "POST":
		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var hello helloJSON
		json.Unmarshal(buf.Bytes(), &hello)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "POST hello! %s\n", hello)

	default:
		fmt.Fprint(w, "Method not allowed.!\n")
	}
}

func Handler2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GET2 hello!\n")
	case "POST":
		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var hello helloJSON
		json.Unmarshal(buf.Bytes(), &hello)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "POST2 hello! %s\n", hello)

	default:
		fmt.Fprint(w, "Method2 not allowed.!\n")
	}
}

func DBInitialize() {
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := "mongodb://root:example@mongo-db:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "あ、どうもお久しぶりです"},
		{Key: "author", Value: "GReeeeN"},
	})
	if err != nil {
		log.Fatal(err)
	}

	episodesResult, err := episodesCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "The Polyglot Developer Podcast"},
		{Key: "author", Value: "Nic Raboy"},
		{Key: "tags", Value: bson.A{"development", "programming", "coding"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(podcastResult, episodesResult)

	cursor, err := podcastsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var podcasts []bson.M
	if err = cursor.All(ctx, &podcasts); err != nil {
		log.Fatal(err)
	}
	fmt.Println(podcasts)
}

func main() {
	DBInitialize()

	router := http.NewServeMux()
	router.HandleFunc("/hello", Handler)
	router.HandleFunc("/pao", Handler2)
	fmt.Printf("[START] server. port: %s\n", addr)
	if err := http.ListenAndServe(addr, middleware.Log(router)); err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
