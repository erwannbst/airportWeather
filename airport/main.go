package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"time"
)

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
	}
	if err := token.Error(); err != nil {

		log.Fatal(err)

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
	IdAirport   int
	MeasureType string
	Value       float32
	Time        string
}

func main() {

	client := connect("tcp://localhost:1883", "airport")

	for true {
		data := &AirportInfo{
			Id:          1,
			IdAirport:   1,
			MeasureType: "Temperature",
			Value:       rand.Float32() * 100,
			Time:        time.Now().String(),
		}

		dataJson, err := json.Marshal(data)

		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		fmt.Println(string(dataJson))

		client.Publish("test", 0, false, string(dataJson))
		time.Sleep(10 * time.Second)
	}
}
