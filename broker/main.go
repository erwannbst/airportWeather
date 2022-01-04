package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

func main() {

	client := connect("tcp://localhost:1883", "my-client-id")

	client.Subscribe("topic", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("%s\n", msg.Payload())
	})
	/*for i := 0; i < 100; i++ {
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