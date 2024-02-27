package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/niyioo/Advanced-Employment-Management-System/backend/models"
)

var client *mongo.Client

func InitMongoDB() error {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the MongoDB server
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func InsertDepartment(department models.Department) error {
    // Get MongoDB collection for departments
    collection := client.Database("AEMS").Collection("departments")

    // Define a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Insert the department into the collection
    _, err := collection.InsertOne(ctx, department)
    if err != nil {
        return err
    }

    return nil
}

// CloseMongoDB closes the connection to MongoDB
func CloseMongoDB() error {
	if client != nil {
		err := client.Disconnect(context.Background())
		return err
	}
	return nil
}
