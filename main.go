package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Todo struct {
    ID        primitive.ObjectID `json:"id" bson:"_id"`
    Completed bool               `json:"completed"`
    Body      string             `json:"body"`
}

var collection *mongo.Collection



func main() {
	fmt.Println("Hello")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(),clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())


	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")
	collection = client.Database("golang_db").Collection("todos")
	app := fiber.New()
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))
	
}
func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return  err
	}

	defer cursor.Close(context.Background())

	for  cursor.Next(context.Background()){
		var todo Todo 
		if err := cursor.Decode(&todo); err != nil {
			return err 
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)

}



// func createTodos(c *fiber.Ctx) error {}
// func updateTodos(c *fiber.Ctx) error {}
// func deleteTodos(c *fiber.Ctx) error {}