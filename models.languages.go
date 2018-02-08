package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	Stats     *LocationStats
}

type LocationStats struct {
	ItemCounts map[string]interface{}
}

type Location struct {
	Languages map[string]language
}

var locationMaps = make(map[string]Location)

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
	currentMap := new(Location)
	currentMap.Languages = make(map[string]language)

	for _, i := range result.Items {
		item := language{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		//change next line to a byte buffer later
		item.Location = officeID
		languageList[itemCounter] = item
		item.Stats = new(LocationStats)
		item.Stats.ItemCounts = make(map[string]interface{})
		item.Stats.ItemCounts["Books"] = item.Books
		item.Stats.ItemCounts["Magazines"] = item.Magazines
		item.Stats.ItemCounts["Bibles"] = item.Bibles
		currentMap.Languages[item.Language+item.Location] = item
		fmt.Println(item.Stats.ItemCounts)
		itemCounter++
	}

	fmt.Println("current map is ", *currentMap)
	locationMaps[officeID] = *currentMap
	return languageList
}

func MakeInventoryChanges(valuemap map[string]interface{}) {
	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		updateLanguage(location, language, itemType, value.(string))
	}

}

func AddInventoryChanges(valuemap map[string]interface{}) {
	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		originalLocation := locationMaps[location].Languages[language+location]
		fmt.Println(originalLocation.Stats.ItemCounts)
		originalValue := originalLocation.Stats.ItemCounts[itemType]

		// originalIntValue, err := strconv.Atoi(originalValue)
		// if err != nil {
		// 	fmt.Println("Query API call failed:")
		// 	fmt.Println((err.Error()))
		// }

		newIntValue, err := strconv.Atoi(value.(string))
		if err != nil {
			fmt.Println("Query API call failed:")
			fmt.Println((err.Error()))

		}
		fmt.Println("location is ", location, " language is ", language,
			" itemType ", itemType, " new value is ", strconv.Itoa(originalValue.(int)+newIntValue))
		updateLanguage(location, language, itemType, strconv.Itoa(originalValue.(int)+newIntValue))
	}

}

func updateLanguage(location string, language string, itemType string, count string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	svc := dynamodb.New(sess)

	// Create item in table SFMetroOfficeInventory
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			fmt.Sprintf(":%s", itemType): {
				N: aws.String(count),
			},
		},
		TableName: aws.String("SFMetroOfficeInventory"),
		Key: map[string]*dynamodb.AttributeValue{
			"Location": {
				S: aws.String(location),
			},
			"Language": {
				S: aws.String(language),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String(fmt.Sprintf("set %s = :%s", itemType, itemType)),
	}

	_, err = svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
