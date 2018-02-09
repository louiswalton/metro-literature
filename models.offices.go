package main

import (
// "fmt"
// "os"

// "github.com/aws/aws-sdk-go/aws"
// "github.com/aws/aws-sdk-go/aws/session"
// "github.com/aws/aws-sdk-go/service/dynamodb"
// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// "github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type office struct {
	ID string `json:"id"`
}

var officeList = []office{
	office{ID: "Depot"},
	office{ID: "Flood Building"},
	office{ID: "Embarcadero"},
	office{ID: "California St"},
	office{ID: "Kearny St."},
}

func getAllOffices() []office {
	return officeList
}
