package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func insertSession(s session) {
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
	sessionDatabase := client.Database("userDatabase")
	sessionCollection := sessionDatabase.Collection("sessionCollection")
	var se session = s

	sessionCollection.InsertOne(ctx, primitive.D{
		{"UUID", se.UUID},
		{"Username", se.Username},
	})
	return
}
func sessionCheck(SessionID string) (session, bool) {
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
	sessionDatabase := client.Database("userDatabase")
	sessionCollection := sessionDatabase.Collection("sessionCollection")
	filter := primitive.D{{"UUID", SessionID}}
	var checkedSession session
	err = sessionCollection.FindOne(ctx, filter).Decode(&checkedSession)
	checkedID := checkedSession.UUID
	var check bool
	if checkedID != "" {
		check = true
	} else {
		check = false
	}
	return checkedSession, check

}

/*
func sessionUpdate(uuid string) {
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
	sessionDatabase := client.Database("userDatabase")
	sessionCollection := sessionDatabase.Collection("sessionCollection")

	filter := bson.M{"uuid": bson.M{"$eq": uuid}}
	update := bson.M{"$set": bson.M{"lastActivity": 42}}

	sessionCollection.UpdateOne(ctx, filter, update)

}
*/
func getUserFromUsername(username string) user {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://darkness:darkness@darkness.30isp.mongodb.net/darkness?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	userDatabase := client.Database("userDatabase")
	usersCollection := userDatabase.Collection("userCollection")

	var usersMatched []user
	var u user
	finder, err := usersCollection.Find(ctx, bson.M{"username": username})
	if err != nil {
		panic(err)
	}
	if err = finder.All(ctx, &usersMatched); err != nil {
		panic(err)
	}
	u = user{usersMatched[0].UserName, usersMatched[0].Password, usersMatched[0].FirstName, usersMatched[0].Last}
	fmt.Println(u.UserName)
	return u

}
