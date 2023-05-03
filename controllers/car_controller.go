package controllers

import (
	"car-dealer-api/configs"
	"car-dealer-api/models"
	"car-dealer-api/responses"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var carCollection *mongo.Collection = configs.GetCollection(configs.DB, "cars")
var validate = validator.New()

// Create a new Car
func CreateCar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var car models.Car
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&car); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CarResponse{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&car); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CarResponse{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newCar := models.Car{
		Id:           primitive.NewObjectID(),
		Model:        car.Model,
		Manufacturer: car.Manufacturer,
		Year:         car.Year,
		Kilometres:   car.Kilometres,
	}

	result, err := carCollection.InsertOne(ctx, newCar)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error: the Car creation process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.CarResponse{Status: http.StatusCreated, Message: "A new Car was added successfully.", Data: &fiber.Map{"data": result}})
}

// Get a Car
func GetCar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	carId := c.Params("carId")
	var car models.Car
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(carId)

	//validate if the car ID exists
	err := carCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&car)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error: invalid car ID.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CarResponse{Status: http.StatusOK, Message: "The operation was successfully.", Data: &fiber.Map{"data": car}})
}

// Edit a Car
func EditCar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	carId := c.Params("carId")
	var car models.Car
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(carId)

	//validate the request body
	if err := c.BodyParser(&car); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CarResponse{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&car); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CarResponse{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"model": car.Model, "manufacturer": car.Manufacturer, "year": car.Year, "kilometres": car.Kilometres}

	result, err := carCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error: the Car edit process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedCar models.Car
	if result.MatchedCount == 1 {
		err := carCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedCar)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error: the Car edit process failed.", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.CarResponse{Status: http.StatusOK, Message: "The Car with the ID " + carId + " was edited correctly.", Data: &fiber.Map{"data": updatedCar}})
}

// Delete a Car
func DeleteCar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	carId := c.Params("carId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(carId)

	result, err := carCollection.DeleteOne(ctx, bson.M{"id": objId})

	//validate if the DeleteOne functions returns an Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error: There is no car with that ID. ", Data: &fiber.Map{"data": err.Error()}})
	}

	//validate the ID number
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.CarResponse{Status: http.StatusNotFound, Message: "Error", Data: &fiber.Map{"data": "Error: The Car with the ID " + carId + " does not exists."}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CarResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": "The Car was deleted successfully."}},
	)
}

// Get All Cars
func GetAllCars(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var cars []models.Car
	defer cancel()

	results, err := carCollection.Find(ctx, bson.M{})

	//validate if the context has a collection
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleCar models.Car
		if err = results.Decode(&singleCar); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CarResponse{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}

		cars = append(cars, singleCar)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CarResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": cars}},
	)
}
