package controllers

import (
	"api/configs"
	"api/models"
	"api/responses"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var sensorCollection *mongo.Collection = configs.GetCollection(configs.DB, "sensor")

func GetAMeasure(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	measureType := c.Params("measure_type")
	//from := c.Params("from")
	//to := c.Params("to")
	var measure models.AirportInfo
	defer cancel()

	err := sensorCollection.FindOne(ctx, bson.M{"MeasureType": measureType}).Decode(&measure)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AirportWeatherResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.AirportWeatherResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": measure}})
}

func GetAllMeasures(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var measures []models.AirportInfo
	defer cancel()


	//matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
	//groupStage := bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}}
	results, err := sensorCollection.Aggregate(ctx, mongo.Pipeline{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AirportWeatherResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleMeasure models.AirportInfo
		if err = results.Decode(&singleMeasure); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AirportWeatherResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		measures = append(measures, singleMeasure)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AirportWeatherResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": measures}})
}