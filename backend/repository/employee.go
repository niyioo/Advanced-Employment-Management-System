package repository

import (
	"context"

	"github.com/niyioo/Advanced-Employment-Management-System/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmployeeRepository struct holds the MongoDB collection reference
type EmployeeRepository struct {
	Collection *mongo.Collection
}

// NewEmployeeRepository creates a new instance of EmployeeRepository
func NewEmployeeRepository(collection *mongo.Collection) *EmployeeRepository {
	return &EmployeeRepository{
		Collection: collection,
	}
}

// CreateEmployee inserts a new employee record into the database
func (r *EmployeeRepository) CreateEmployee(ctx context.Context, employee *models.Employee) error {
	_, err := r.Collection.InsertOne(ctx, employee)
	return err
}

// GetEmployeeByID retrieves an employee record from the database by ID
func (r *EmployeeRepository) GetEmployeeByID(ctx context.Context, id string) (*models.Employee, error) {
	var employee models.Employee
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// UpdateEmployee updates an existing employee record in the database
func (r *EmployeeRepository) UpdateEmployee(ctx context.Context, id string, employee *models.Employee) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": employee})
	return err
}

// DeleteEmployee deletes an employee record from the database by ID
func (r *EmployeeRepository) DeleteEmployee(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetAllEmployees retrieves all employee records from the database
func (r *EmployeeRepository) GetAllEmployees(ctx context.Context) ([]models.Employee, error) {
	var employees []models.Employee
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &employees)
	if err != nil {
		return nil, err
	}
	return employees, nil
}
