package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type product struct {
	ProductName string
	Condition   string
	Rentable    bool
	Price       string
	Description string
}

func castBool(s string) bool {
	b := true
	f := "false"
	if s == f {
		b = false
	}
	return b
}

func insertProduct(p product) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://darkness:darkness@darkness.30isp.mongodb.net/darkness?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	productDatabase := client.Database("productDatabase")
	productCollection := productDatabase.Collection("productCollection")

	productCollection.InsertOne(ctx, primitive.D{
		{"ProductName", p.ProductName},
		{"Condition", p.Condition},
		{"Rentable", p.Rentable},
		{"Price", p.Price},
		{"Description", p.Description},
	})
	return
}

func getAllProducts() []product {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://darkness:darkness@darkness.30isp.mongodb.net/darkness?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	productDatabase := client.Database("productDatabase")
	productCollection := productDatabase.Collection("productCollection")

	var productFind []product

	finder, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	if err = finder.All(ctx, &productFind); err != nil {
		panic(err)
	}
	return productFind
}
