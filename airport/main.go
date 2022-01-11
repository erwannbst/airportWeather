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

// Config structure get from config.json
type Config struct {
	BrokerAddress string
	BrokerPort    int
	QoS           byte
	ClientId      string // equals to airport IATA code
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

// Airport data send to the broker
type AirportInfo struct {
	IdSensor    int
	IdAirport   string
	MeasureType string
	Value       float32
	Time        string
}

/*
	Create a sensor with provided data and publish one time to the broker
	params :
		- client : broker client
		- airportId : IATA code of the airport
		- sensorId : id of the sensor
		- measureType : the name of the thing we measure
		- qos : level of quality of service
		- generateValue : function to calculate the value of the sensor
*/
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

	client.Publish("airport", qos, false, string(dataJson))
}

// Generate a value for temperature between 0 and 100
func generateTemp() float32 {
	return rand.Float32() * 100
}

// Generate a value for atmospheric pressure between 990 and 1050
func generatePress() float32 {
	return 990 + rand.Float32()*60
}

// Generate a value for wind speed between 0 and 120
func generateWind() float32 {
	return rand.Float32() * 120
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

	// Connect to the broker
	client := connect(config.BrokerAddress+":"+strconv.Itoa(config.BrokerPort), config.ClientId)

	// Create sensors and send data to the broker every 10 seconds
	for true {
		createSensor(client, config.ClientId, 1, "Temperature", config.QoS, generateTemp)
		createSensor(client, config.ClientId, 2, "Atmospheric pressure", config.QoS, generatePress)
		createSensor(client, config.ClientId, 3, "Wind speed", config.QoS, generateWind)
		time.Sleep(10 * time.Second)
	}
}
