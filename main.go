package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}


type Post struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Caption   string    `json:"caption,omitempty" bson:"caption,omitempty"`
	Imageurl  string    `json:"imgurl,omitempty" bson:"imgurl,omitempty"`
	UserId    string    `json:"uid,omitempty" bson:"uid,omitempty"`
	Timestamp time.Time `json:"created,omitempty" bson:"created,omitempty"`
}

func main() {
	connect()
	handleRequest()
}

var client *mongo.Client

func connect() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false")
	client, _ = mongo.NewClient(clientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("error", err)
	} else {
		log.Println("Connected to MondoDB Server")
	}

}

func handleRequest() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", returnAllUsers)
	http.HandleFunc("/users/", returnSingleUser)
	http.HandleFunc("/posts", returnAllPosts)
	http.HandleFunc("/posts/", returnSinglePost)
	http.HandleFunc("/posts/users/", returnAllUserPosts)
	err := http.ListenAndServe(":6000", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Khyati Heda")
	fmt.Println(" Home Page")
}

func returnAllUsers(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var users []User
		collection := client.Database("golang").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: returnAllUsers")
		json.NewEncoder(response).Encode(users)
	} else {
		request.ParseForm()
		decoder := json.NewDecoder(request.Body)
		var newUser User
		err := decoder.Decode(&newUser)
		if err != nil {
			panic(err)
		}
		log.Println(newUser.ID)
		fmt.Println("User Added")
		insertUser(newUser)
	}
}

func returnSingleUser(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	var id string = request.URL.Path
	id = id[7:]
	var user User
	collection := client.Database("golang").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Println("Returned User ID NO : ", user.ID)
	json.NewEncoder(response).Encode(user)
}

func insertUser(user User) {
	collection := client.Database("golang").Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted user with ID:", insertResult.InsertedID)
}

func returnAllPosts(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var posts []Post
		collection := client.Database("golang").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println(" returnAllPosts")
		json.NewEncoder(response).Encode(posts)
	} else {
		request.ParseForm()
		decoder := json.NewDecoder(request.Body)
		var newPost Post
		newPost.Timestamp = time.Now()
		err := decoder.Decode(&newPost)
		if err != nil {
			panic(err)
		}
		log.Println(newPost.ID)
		fmt.Println("Post Added")
		insertPost(newPost)
	}
}

func returnSinglePost(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	var id string = request.URL.Path
	id = id[7:]
	var posst Post
	collection := client.Database("golang").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, Post{ID: id}).Decode(&posst)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Println("Returned Post ID NO : ", posst.ID)
	json.NewEncoder(response).Encode(posst)
}

func insertPost(post Post) {
	collection := client.Database("golang").Collection("posts")
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

func returnAllUserPosts(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	var id string = request.URL.Path
	id = id[13:]
	collection := client.Database("golang").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, Post{UserId: id})
	for cursor.Next(ctx) {

		var episodesFiltered []bson.M
		if err != nil {
			log.Fatal(err)
		}
		if err = cursor.All(ctx, &episodesFiltered); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(response).Encode(episodesFiltered)
	}

}