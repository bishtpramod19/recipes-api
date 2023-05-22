// Recipes API
//
// This is a sample recipes API

//  Schemes: http
//  Host: localhost:8080
//  BasePath: /
//  Version: 1.0.0
//  Contact: Pramod Singh Bisht
// <bisht.pramod19@gmail.com>
//

//  Consumes:
//  - application/json
//

//  Produces:
//  - application/json
// swagger:meta

package main

import (
	"context"
	"log"
	"os"

	"github.com/bishtpramod19/recipes-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DATABASE", "RecipesData")
	log.Println("environment variable set now")
	// recipes = make([]Recipe, 0)

	// file, _ := ioutil.ReadFile("recipes.json")
	// err = json.Unmarshal([]byte(file), &recipes)

	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	//redis initialization

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// redisStatus := redisClient.Ping()
	// fmt.Println(redisStatus)

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("Recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.Run()
}
