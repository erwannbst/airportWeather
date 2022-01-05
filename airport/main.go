package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BrokerAddress string
	BrokerPort    int
	QoS           byte
	ClientId      string
}

func createClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)
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
	IdSensor    int
	IdAirport   string
	MeasureType string
	Value       float32
	Time        string
}

func createSensor(client mqtt.Client, airportId string, sensorId int, measureType string, qos byte, generateValue func() float32) {

	data := &AirportInfo{
		IdSensor:    sensorId,
		IdAirport:   airportId,
		MeasureType: measureType,
		Value:       generateValue(),
		Time:        time.Now().Format("2006-01-02-15-04-05"),
	}

	dataJson, err := json.Marshal(data)

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Println(string(dataJson))

	client.Publish("test", qos, false, string(dataJson))
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

	// Read config file
	var config Config
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	client := connect(config.BrokerAddress+":"+strconv.Itoa(config.BrokerPort), config.ClientId)

	for true {
		createSensor(client, config.ClientId, 1, "Temperature", config.QoS, generateTemp)
		createSensor(client, config.ClientId, 2, "Atmospheric pressure", config.QoS, generatePress)
		createSensor(client, config.ClientId, 3, "Wind speed", config.QoS, generateWind)
		time.Sleep(10 * time.Second)
	}
}
