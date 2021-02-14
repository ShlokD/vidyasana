package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Course struct {
	Title string `bson:title`
}

func connectToDb() (*mongo.Client, context.Context) {
	mongoURI := os.Getenv("MONGODB_ATLAS_URI")
	atlasURI := fmt.Sprintf("%s?retryWrites=true&w=majority", mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(atlasURI))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}
func handleIndex(c *gin.Context) {
	client, ctx := connectToDb()
	defer client.Disconnect(ctx)

	dbName := os.Getenv("MONGODB_DB_NAME")
	var course Course
	col := client.Database(dbName).Collection("courses")
	err := col.FindOne(context.TODO(), bson.D{}).Decode(&course)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
	}
	c.JSON(http.StatusOK, gin.H{"title": course.Title})
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", handleIndex)

	router.Run(":" + port)
}
