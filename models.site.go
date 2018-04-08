package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type site struct {
	Location string `json:"id"`
	Office   string `json:"office"`
}

type LocationReport struct {
	Report          map[string]string `json:"Report"`
	ReportTimestamp string            `json:"ReportTimestamp"`
	ReportID        string            `json:"ReportID"`
	Location        string            `json:"Location"`
	Office          string            `json:"Office"`
	ReportDate      string            `json:"ReportDate"`
}

func getSiteByOfficeID(officeID string) []site {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	filt := expression.Name("Office").Equal(expression.Value(officeID))
	proj := expression.NamesList(expression.Name("ID"))

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
		TableName:                 aws.String("Site"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	siteList := make([]site, len(result.Items))

	for j, i := range result.Items {
		item := site{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		siteList[j] = item
	}
	//fmt.Println(siteList)
	return siteList
}

func CreateReport(report map[string]string, location string, office string,
	date string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	var Pacific = "America/Los_Angeles"
	loc, _ := time.LoadLocation(Pacific)
	var t = time.Now().In(loc)
	fmt.Println(t.String())
	fmt.Println(date)
	newReport := LocationReport{
		Report:          report,
		ReportID:        t.String() + location,
		ReportTimestamp: t.String(),
		Location:        location,
		Office:          office,
		ReportDate:      date,
	}

	av, err := dynamodbattribute.MarshalMap(newReport)

	// Create item in table LocationReports
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("LocationReports"),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func GetLocationReports(location string) []LocationReport {
	fmt.Println(location)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	filt := expression.Name("Location").Equal(expression.Value(location))
	proj := expression.NamesList(
		expression.Name("ReportID"), expression.Name("ReportTimestamp"),
		expression.Name("Location"), expression.Name("Office"),
		expression.Name("Report"), expression.Name("ReportDate"))

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
		TableName:                 aws.String("LocationReports"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	reportList := make([]LocationReport, len(result.Items))

	for j, i := range result.Items {
		item := LocationReport{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		reportList[j] = item
	}

	return reportList
}
