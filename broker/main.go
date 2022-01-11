package main

import (
	"context"
	"encoding/csv"
	_ "encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// AirportInfo Data structure of the incoming messages/**
type AirportInfo struct {
	IdSensor    int
	IdAirport   string
	MeasureType string
	Value       float32
	Time        string
}

/**
listen the publishers
 */
func main() {

	//connection to the broker
	client := connect("tcp://localhost:1883", "my-client-id")


	var wg sync.WaitGroup
	wg.Add(1)
	//save the messages received into database and csv file
	go func() {
		subscribe(client, "airport", 0, func(client mqtt.Client, msg mqtt.Message) {
			saveFile(string(msg.Payload()))
			saveDb(string(msg.Payload()))
		})
	}()
	wg.Wait()

}

/**
Save the message into Mongodb Atlas
 */
func saveDb(message string) {

	//convert message into json format
	var info AirportInfo
	json.Unmarshal([]byte(message), &info)

	//connect to the db
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://Mael:Argenttropbien@cluster0.5j16q.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	check(err)

	//move to the collection sensor in the airport database
	airportDatabase := client.Database("airport")
	sensorCollection := airportDatabase.Collection("sensor")

	//insert data
	_, err = sensorCollection.InsertOne(ctx, bson.D{
		{Key: "IdSensor", Value: info.IdSensor},
		{Key: "IdAirport", Value: info.IdAirport},
		{Key: "MeasureType", Value: info.MeasureType},
		{Key: "Value", Value: info.Value},
		{Key: "Time", Value: info.Time},
	})

	check(err)

	fmt.Println("Inserted data into sensor collection !")
}

/**
Save message in a csv file
 */
func saveFile(message string) {

	//convert message into json format
	var info AirportInfo
	json.Unmarshal([]byte(message), &info)

	//path of the datalake
	filename := "../datalake/"

	t, err := time.Parse("2006-01-02-15-04-05", info.Time)
	check(err)

	//format the date
	var day string
	if t.Day() < 10 {
		day = "0" + strconv.Itoa(t.Day())
	} else {
		day = strconv.Itoa(t.Day())
	}

	var month string
	if int(t.Month()) < 10 {
		month = "0" + strconv.Itoa(int(t.Month()))
	} else {
		month = strconv.Itoa(int(t.Month()))
	}

	year := strconv.Itoa(t.Year())

	//name of the csv file
	filename += info.IdAirport + "-" + year + "-" + month + "-" + day + ".csv"
	var f *os.File

	_, err = os.Stat(filename)
	fileNotExists := errors.Is(err, os.ErrNotExist)

	//create/open the file
	f, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	check(err)

	//close the file at the end of the block to avoid error when opening the file during execution
	defer f.Close()

	//If the file not exists set the separator and create the columns
	if fileNotExists {
		_, err = f.WriteString("sep=,\n\"IdSensor\", \"IdAirport\", \"MeasureType\", \"Value\", \"Time\"\n")
		check(err)

	}

	csvwriter := csv.NewWriter(f)

	//set the message
	empData := [][]string{

		{strconv.Itoa(info.IdSensor), info.IdAirport, info.MeasureType, fmt.Sprintf("%f", info.Value), info.Time},
	}

	//write the message
	for _, empRow := range empData {
		_ = csvwriter.Write(empRow)
	}

	csvwriter.Flush()

	check(err)
	fmt.Println("Inserted data into " + filename + " !")

}

/**
Manage the errors
 */
func check(e error) {
	if e != nil {
		panic(e)
	}
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

func subscribe(client mqtt.Client, topic string, qos byte, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, qos, callback)
	token.Wait()
}
