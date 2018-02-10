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

type StockItem struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type Stocklist struct {
	Stocklist      map[string]StockItem
	HasBeenChanged bool
}

var Stock = Stocklist{Stocklist: make(map[string]StockItem), HasBeenChanged: false}

func getStocklist() Stocklist {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	proj := expression.NamesList(expression.Name("Language"), expression.Name("Code"))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Stock"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	var newStockList = Stocklist{Stocklist: make(map[string]StockItem), HasBeenChanged: false}

	for _, i := range result.Items {
		item := StockItem{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		//change next line to a byte buffer later
		newStockList.Stocklist[item.Language] = item
	}

	Stock = newStockList
	return Stock
}
