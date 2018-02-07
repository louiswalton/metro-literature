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

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-west-1")},
	// )

	// if err != nil {
	// 	fmt.Println("Got error creating session:")
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// svc := dynamodb.New(sess)

	// filt := expression.Name("Location").AttributeExists()
	// // Get back the Location
	// proj := expression.NamesList(expression.Name("Location"))

	// expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	// if err != nil {
	// 	fmt.Println("Got error building expression:")
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// // Build the query input parameters
	// params := &dynamodb.ScanInput{
	// 	ExpressionAttributeNames:  expr.Names(),
	// 	ExpressionAttributeValues: expr.Values(),
	// 	FilterExpression:          expr.Filter(),
	// 	ProjectionExpression:      expr.Projection(),
	// 	TableName:                 aws.String("Locations"),
	// }

	// // Make the DynamoDB Query API call
	// result, err := svc.Scan(params)

	// if err != nil {
	// 	fmt.Println("Query API call failed:")
	// 	fmt.Println((err.Error()))
	// 	os.Exit(1)
	// }

	// itemCounter := 0
	// officeList := make([]office, *result.ScannedCount)

	// for _, i := range result.Items {
	// 	item := office{}

	// 	err = dynamodbattribute.UnmarshalMap(i, &item)

	// 	if err != nil {
	// 		fmt.Println("Got error unmarshalling:")
	// 		fmt.Println(err.Error())
	// 		os.Exit(1)
	// 	}
	// 	officeList[itemCounter] = item
	// 	itemCounter++

	// }
	// fmt.Println(officeList)
	return officeList
}
