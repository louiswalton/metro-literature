package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type language struct {
	Location  string `json:"location"`
	Language  string `json:"language"`
	Books     int    `json:"books"`
	Magazines int    `json:"magazines"`
	Bibles    int    `json:"bibles"`
}

func getAllLanguages() []language {

	return getInventoryByOfficeID("Depot")
}

func getInventoryByOfficeID(officeID string) []language {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	filt := expression.Name("Location").Equal(expression.Value(officeID))
	proj := expression.NamesList(expression.Name("Language"), expression.Name("Bibles"),
		expression.Name("Books"), expression.Name("Magazines"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("SFMetroOfficeInventory"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	languageList := make([]language, *result.Count)
	itemCounter := 0

	for _, i := range result.Items {
		item := language{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		languageList[itemCounter] = item
		itemCounter++
	}
	return languageList
}
