package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8084"
	rpcPort  = "5001"
	rRpcPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
}

func main() {
	clientMongo, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}
	client = clientMongo

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = clientMongo.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

}

func connectToMongo() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Printf("Error connecting to mongo %s", err)
		return nil, err
	}

	return c, nil

}
