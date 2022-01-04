package main

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	client := connect("tcp://localhost:1883", "my-client-id")

	wg.Add(1)
	go func() {
		subscribe(client, "test", 0, func(client mqtt.Client, msg mqtt.Message) {
			saveFile(string(msg.Payload()))
			saveDb(string(msg.Payload()))
		})
	}()

	wg.Wait()
	/*
		for i := 0; i < 100; i++ {
			client.Publish("topic", 0, false, "emilien est nul")
		}*/

}

func createClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	// AddBroker adds a broker URI to the list of brokers to be used.
	// The format should be "scheme://host:port"
	opts.AddBroker(brokerURI)
	// opts.SetUsername(user)
	// opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts

}

func connect(brokerURI string, clientId string) mqtt.Client {

	fmt.Println("Trying to connect (" + brokerURI + ", " + clientId + ")...")
	opts := createClientOptions(brokerURI, clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
		if err := token.Error(); err != nil {

			log.Fatal(err)

		}
	}
	return client

}

func publish(client mqtt.Client, topic string, qos byte, payload string) {
	token := client.Publish(topic, qos, false, payload)
	token.Wait()
}

func subscribe(client mqtt.Client, topic string, qos byte, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, qos, callback)
	token.Wait()
}

func unsubscribe(client mqtt.Client, topic string) {
	token := client.Unsubscribe(topic)
	token.Wait()
}

func disconnect(client mqtt.Client) {
	client.Disconnect(250)
}

type AirportInfo struct {
	Id          int
	IdAirport   string
	MeasureType string
	Value       float32
	Time        string
}

func saveDb(message string) {
	var info AirportInfo
	json.Unmarshal([]byte(message), &info)
	t, error := time.Parse("yyyy-mm-dd-hh-mm-ss", info.Time)
	var date string
	if error != nil {
		date = t.String()
	}



	clientOptions := options.Client().
		ApplyURI("mongodb+srv://Mael:Argenttropbien@cluster0.5j16q.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	airportDatabase := client.Database("airport")
	sensorCollection := airportDatabase.Collection("sensor")

	_, err = sensorCollection.InsertOne(ctx, bson.D{
		{Key: "IdSensor", Value: info.Id},
		{Key: "IdAirport", Value: info.IdAirport},
		{Key: "MeasureType", Value: info.MeasureType},
		{Key: "Value", Value: info.Value},
		{Key: "Time", Value: date},
	})

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Inserted documents into sensor collection !")
}

func saveFile(message string) {
	var info AirportInfo
	json.Unmarshal([]byte(message), &info)
	filename := "C:\\Users\\maels\\Documents\\imt\\archiD\\Go\\airportWeather\\"
	filename += info.IdAirport
	t, error := time.Parse("yyyy-mm-dd-hh-mm-ss", info.Time)
	if error != nil {
		filename += strconv.Itoa(t.Day()) +t.Month().String() + strconv.Itoa(t.Year())
	}
	filename += ".csv"
	fmt.Println(filename)

}
