package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"time"
)

var AIRPORT_ID = "CDG"

func createClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
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
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
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

func createSensor(client mqtt.Client, sensorId int, name string, generateValue func() float32) {

	data := &AirportInfo{
		Id:          sensorId,
		IdAirport:   AIRPORT_ID,
		MeasureType: name,
		Value:       generateValue(),
		Time:        time.Now().String(),
	}

	dataJson, err := json.Marshal(data)

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Println(string(dataJson))

	client.Publish("test", 0, false, string(dataJson))
}

func generateTemp() float32 {
	return rand.Float32() * 100
}

func generatePress() float32 {
	return 990 + rand.Float32()*60
}

func generateWind() float32 {
	return rand.Float32() * 100
}

func main() {

	client := connect("tcp://localhost:1883", AIRPORT_ID)

	for true {
		createSensor(client, 1, "Temperature", generateTemp)
		createSensor(client, 2, "Atmospheric pressure", generatePress)
		createSensor(client, 3, "Wind speed", generateWind)
		time.Sleep(10 * time.Second)
	}
}
