package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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

	currentMap := new(Location)
	currentMap.Languages = make(map[string]language)
	currentMap.Location = officeID

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
		item.Stats = new(LocationStatsMap)
		item.Stats.ItemCounts = make(map[string]interface{})
		item.Stats.ItemCounts["Books"] = item.Books
		item.Stats.ItemCounts["Magazines"] = item.Magazines
		item.Stats.ItemCounts["Bibles"] = item.Bibles
		currentMap.Languages[item.Language+item.Location] = item
	}

	locationMaps[officeID] = processLanguageList(*currentMap, getStocklist())

	languageList := make([]language, len(locationMaps[officeID].Languages))
	itemCounter := 0
	for _, v := range locationMaps[officeID].Languages {
		languageList[itemCounter] = v
		itemCounter++
	}
	sort.Sort(sort.Reverse(lList(languageList)))
	return languageList
}

func processLanguageList(officeInventory Location, stock Stocklist) Location {
	for _, item := range stock.Stocklist {
		_, ok := officeInventory.Languages[item.Language+officeInventory.Location]
		if !ok {
			newItem := language{Language: item.Language, Location: officeInventory.Location,
				Books: 0, Magazines: 0, Bibles: 0}
			newItem.Stats = new(LocationStatsMap)
			newItem.Stats.ItemCounts = make(map[string]interface{})
			newItem.Stats.ItemCounts["Books"] = 0
			newItem.Stats.ItemCounts["Magazines"] = 0
			newItem.Stats.ItemCounts["Bibles"] = 0
			officeInventory.Languages[item.Language+officeInventory.Location] = newItem
		}
	}
	return officeInventory
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

func MakeInventoryChanges(valuemap map[string]interface{}) {
	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		if len(value.(string)) != 0 {
			updateLanguage(location, language, itemType, value.(string))
		}
	}

}

func AddInventoryChanges(valuemap map[string]interface{}) {
	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		originalLocation := locationMaps[location].Languages[language+location]
		originalValue := originalLocation.Stats.ItemCounts[itemType]
		if len(value.(string)) != 0 {
			newIntValue, err := strconv.Atoi(value.(string))
			if err != nil {
				fmt.Println("Query API call failed:")
				fmt.Println((err.Error()))
			}

			updateLanguage(location, language, itemType, strconv.Itoa(originalValue.(int)+newIntValue))
		}
	}
}

func SubtractInventoryFromDepot(valuemap map[string]interface{}) {
	getInventoryByOfficeID("Depot")
	originalLocation := locationMaps["Depot"]
	originalLanguages := originalLocation.Languages

	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		if len(value.(string)) != 0 {
			originalValue := originalLanguages[language+"Depot"].Stats.ItemCounts[itemType]
			newIntValue, err := strconv.Atoi(value.(string))
			if err != nil {
				fmt.Println("Query API call failed:")
				fmt.Println((err.Error()))
			}
			if location != "Depot" {
				updateLanguage("Depot", language, itemType, strconv.Itoa(originalValue.(int)-newIntValue))
			}
		}
	}
}
