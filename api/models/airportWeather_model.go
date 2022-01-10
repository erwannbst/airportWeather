package models

// AirportInfo Data structure of the incoming messages/**
type AirportInfo struct {
	IdSensor  		int        `json:"IdSensor"`
	IdAirport 		string     `json:"IdAirport"`
	MeasureType  	string     `json:"MeasureType"`
	Value 			float32     `json:"Value"`
	Time 			string     `json:"Time"`
}