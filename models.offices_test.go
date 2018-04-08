package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

import "testing"

// Test the function that gets all of the languages
func TestLogTransaction(t *testing.T) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	logTransaction(svc, "test_create_location", "test_create_language", "test_create_book", "4", "add")
}

func TestGetTransactionLog(t *testing.T) {
	logs := getInventoryLog("test_create_location")
	fmt.Println(logs)
	fmt.Println(len(logs))
	if len(logs) != 3 {
		t.Fail()
	}
}
