package db

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongo instance type
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

// define MongoInstance variable
var mg MongoInstance

// define constants for MongoInstance
const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017/" + dbName //todo: may need to include username@password

// connect to Mongodb
func Connect() error {
	// create a mongo client, set options and apply for mongoURI
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	// if 30 seconds elapses w/ no successful connection, abort Connect() (i.e. stop trying to connect to mongodb)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// connect to the mongodb server
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	// get handle for database
	db := client.Database(dbName)

	// set up Mongo instance
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

// routes
func GetEmployee(ctx *fiber.Ctx) error {
	// initialize empty slice of employeees
	var employees []Employee = make([]Employee, 0)

	// get all employees in mongodb
	query := bson.D{{}}

	// find all employees
	cursor, err := mg.dB.Collection("employees").Find(ctx.Context(), query)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	// iterate through all employee info, write to employees slice, pass to context
	if err := cursor.All(ctx.Context(), &employees); err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	// return employee info as JSON
	return ctx.JSON(employees)
}

func AddEmployee(ctx *fiber.Ctx) error {
	// initialize employee variable
	employee := new(Employee)

	// get employees collection
	collection := mg.Db.Collection("employees")

	// parse out employee info from body, store in employee
	if err := ctx.BodyParser(employee); err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	// intialize employee ID
	employee.ID = ""

	// insert new employee into db
	insertResult, err := collection.InsertOne(ctx.Context(), employee)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	// find new employee record using insertResult ID
	query := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	newRecord := collection.FindOne(ctx.Context(), query)

	// return created record to client
	newEmployee := &Employee{}
	newRecord.Decode(newEmployee)

	return ctx.Status(201).JSON(newEmployee)
}

func UpdateEmployee(ctx *fiber.Ctx) error {
	//
}

func DeleteEmployee(ctx *fiber.Ctx) error {

}
