package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/models" // Import your models package
)

var client *mongo.Client

// Define Employee struct according to the data model
type Employee struct {
	ID        string `json:"id" bson:"_id"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	// Add more fields as needed
}

func main() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping the MongoDB server
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	// Define InsertDepartment function
	InsertDepartment := func(department models.Department) error {
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

	r := gin.Default()

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Advanced Employment Management System!",
		})
	})

	// Define routes for employee management
	r.GET("/employees", func(c *gin.Context) {
		// Get MongoDB collection
		collection := client.Database("AEMS").Collection("employees")

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find all employees
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		// Iterate over the cursor and collect employees
		var employees []Employee
		for cursor.Next(ctx) {
			var employee Employee
			if err := cursor.Decode(&employee); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			employees = append(employees, employee)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, employees)
	})

	// Define route for fetching a specific employee by ID
	r.GET("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Get MongoDB collection
		collection := client.Database("AEMS").Collection("employees")

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Define filter to find the employee by ID
		filter := bson.M{"_id": id}

		// Find the employee by ID
		var employee Employee
		err := collection.FindOne(ctx, filter).Decode(&employee)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		c.JSON(http.StatusOK, employee)
	})

	// Define route for creating a new employee
	r.POST("/employees", func(c *gin.Context) {
		var employee Employee
		if err := c.BindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get MongoDB collection
		collection := client.Database("AEMS").Collection("employees")

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Insert the new employee into the collection
		_, err := collection.InsertOne(ctx, employee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Employee created successfully"})
	})

	// Define route for updating an existing employee
	r.PUT("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")

		var employee Employee
		if err := c.BindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get MongoDB collection
		collection := client.Database("AEMS").Collection("employees")

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Define filter to find the employee by ID
		filter := bson.M{"_id": id}

		// Define update to set the new employee data
		update := bson.M{"$set": employee}

		// Perform the update operation
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
	})

	// Define route for deleting an existing employee
	r.DELETE("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Get MongoDB collection
		collection := client.Database("AEMS").Collection("employees")

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Define filter to find the employee by ID
		filter := bson.M{"_id": id}

		// Perform the delete operation
		_, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
	})

	// Run the server
	r.Run(":8080")
}
