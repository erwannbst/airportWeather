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
	"os"
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
	IdSensor    int
	IdAirport   string
	MeasureType string
	Value       float32
	Time        string
}

func saveDb(message string) {
	var info AirportInfo
	json.Unmarshal([]byte(message), &info)
	t, err := time.Parse("yyyy-mm-dd-hh-mm-ss", info.Time)
	var date string
	if err != nil {
		date = t.String()
	}



	clientOptions := options.Client().
		ApplyURI("mongodb+srv://Mael:Argenttropbien@cluster0.5j16q.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	check(err)

	airportDatabase := client.Database("airport")
	sensorCollection := airportDatabase.Collection("sensor")

	_, err = sensorCollection.InsertOne(ctx, bson.D{
		{Key: "IdSensor", Value: info.IdSensor},
		{Key: "IdAirport", Value: info.IdAirport},
		{Key: "MeasureType", Value: info.MeasureType},
		{Key: "Value", Value: info.Value},
		{Key: "Time", Value: date},
	})

	check(err)

	fmt.Println("Inserted documents into sensor collection !")
}

func saveFile(message string) {
	var info AirportInfo
	json.Unmarshal([]byte(message), &info)
	filename := "C:\\Users\\maels\\Documents\\imt\\archiD\\Go\\airportWeather\\"
	filename += info.IdAirport
	t, error := time.Parse("2006-01-02-15-04-05", info.Time)

	var day string

	if(t.Day()<10){
		day = "0" + strconv.Itoa(t.Day())
	}else{
		day= strconv.Itoa(t.Day())
	}

	var month string

	if(int(t.Month()) <10){
		month = "0" + strconv.Itoa(int(t.Month()))
	}else{
		month = strconv.Itoa(int(t.Month()))
	}

	year := strconv.Itoa(t.Year())

	if error != nil {
		filename += day +month + year
	}
	filename += ".csv"

	fmt.Println(filename)

	f, err := os.Create(filename)
	check(err)
	_, err = f.WriteString(message)
	check(err)
	fmt.Println("information wrote in file")


}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
