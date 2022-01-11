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

func GetMeasureByType(c *fiber.Ctx) error {

	measureType := c.Query("measure_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var measures []models.AirportInfo
	defer cancel()

	//db.sensor.aggregate([{$match:{"MeasureType":"Wind speed","Time":{$gte:"2022-01-09",$lte:"2022-01-10"}}}])

	matchStage := bson.D{{
		"$match", bson.D{{"MeasureType", measureType},{"Time",bson.D{{"$gte",startDate},{"$lte", endDate}}}}}}

	results ,err := sensorCollection.Aggregate(ctx, mongo.Pipeline{matchStage})
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

func GetAllMeasures(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var measures []models.AverageMeasuresInfo
	defer cancel()



	day := c.Query("day")

	matchStage := bson.D{{
		"$match", bson.D{{"Time", bson.M{"$gte": day+"-00-00-00", "$lte": day+"-23-59-59"}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$MeasureType"}, {"AvgValue", bson.D{{"$avg", "$Value"}}}}}}
	results, err := sensorCollection.Aggregate(ctx, mongo.Pipeline{matchStage,groupStage})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AirportWeatherResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleMeasure models.AverageMeasuresInfo
		if err = results.Decode(&singleMeasure); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AirportWeatherResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		measures = append(measures, singleMeasure)

	}

	return c.Status(http.StatusOK).JSON(
		responses.AirportWeatherResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": measures}})
}