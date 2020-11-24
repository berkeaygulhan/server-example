package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	UserName  string
	Password  []byte
	FirstName string
	Last      string
}

type session struct {
	Username string
	UUID     string
}

func getUser(w http.ResponseWriter, req *http.Request) user {
	// get cookie
	c, err := req.Cookie("session")

	if err != nil {
		sID, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		c = &http.Cookie{
			Name:  "session2",
			Value: sID.String(),
		}

	}
	//c.MaxAge = sessionLenght
	http.SetCookie(w, c)

	// if the user exists already, get user
	var u user
	if s, ok := sessionCheck(c.Value); ok {
		//sessionUpdate(s.UUID)
		fmt.Println(s.UUID)
		u = getUserFromUsername(s.Username)

	}

	return u
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	logged := true
	c, err := req.Cookie("session")
	if err != nil {
		logged = false
	}

	s, ok := sessionCheck(c.Value)
	if ok {
		//s.lastActivity = time.Now()
	}
	if true != checkUsernameUnique(s.Username) {
		logged = false
	}

	//c.MaxAge = sessionLenght
	http.SetCookie(w, c)
	return logged

}

/*
s.lastActivity = time.Now()
//overwrite
user1 := findUsername(s.un)
log.Println(user1)
func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
}
*/
func insertUser(un user) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://darkness:darkness@darkness.30isp.mongodb.net/darkness?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var u user = un

	defer client.Disconnect(ctx)
	userDatabase := client.Database("userDatabase")
	usersCollection := userDatabase.Collection("userCollection")

	userResult, err := usersCollection.InsertOne(ctx, primitive.D{
		{"firstName", u.FirstName},
		{"lastName", u.Last},
		{"username", u.UserName},
		{"password", u.Password},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userResult.InsertedID)
}

func findUser(username string, password string) bool {
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
	userDatabase := client.Database("userDatabase")
	usersCollection := userDatabase.Collection("userCollection")

	var usersMatched []user
	passwordCorrect := true
	finder, err := usersCollection.Find(ctx, bson.M{"username": username})
	if err != nil {
		panic(err)
	}
	if err = finder.All(ctx, &usersMatched); err != nil {
		panic(err)
	}
	userData := usersMatched[0].Password
	err = bcrypt.CompareHashAndPassword(userData, []byte(password))
	if err != nil {
		passwordCorrect = false
	}
	return passwordCorrect

}

func checkUsernameUnique(username string) bool {
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
	userDatabase := client.Database("userDatabase")
	usersCollection := userDatabase.Collection("userCollection")

	var usernameMatched []user
	usernameUnique := false

	finder, err := usersCollection.Find(ctx, bson.M{"username": username})
	if err != nil {
		panic(err)
	}
	if err = finder.All(ctx, &usernameMatched); err != nil {
		panic(err)
	}
	if len(usernameMatched) < 1 {
		usernameUnique = true
	}
	return usernameUnique

}
