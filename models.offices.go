package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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
	office{ID: "Mill Valley"},
	office{ID: "Golden Gate Park"},
}

type transactionLog struct {
	Location        string `json:"Location"`
	Language        string `json:"Language"`
	TransactionID   string `json:"TransactionID"`
	Transaction     string `json:"Transaction"`
	ItemType        string `json:"ItemType"`
	Count           string `json:"Count"`
	TransactionTime string `json:"TransactionTime"`
}

const officeInventoryTableName string = "SanFrancisco-OfficeInventory"
const transactionLogTableName string = "SanFrancisco-TransactionLog"

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
		TableName:                 aws.String(officeInventoryTableName),
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

func updateLanguage(svc *dynamodb.DynamoDB, location string, language string, itemType string, count string) {

	// Create item in table SanFrancisco-OfficeInventory
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			fmt.Sprintf(":%s", itemType): {
				N: aws.String(count),
			},
		},
		TableName: aws.String(officeInventoryTableName),
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

	_, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func logTransaction(svc *dynamodb.DynamoDB, location string, language string, itemType string, count string,
	transaction string) {

	var Pacific = "America/Los_Angeles"
	loc, _ := time.LoadLocation(Pacific)
	var t = time.Now().In(loc)
	// fmt.Println(t.String())

	newTransaction := transactionLog{
		Location:        location,
		Language:        language,
		TransactionID:   t.String() + location,
		Transaction:     transaction,
		ItemType:        itemType,
		Count:           count,
		TransactionTime: t.String(),
	}

	av, err := dynamodbattribute.MarshalMap(newTransaction)

	// Create item in table TransactionLog
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(transactionLogTableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func MakeInventoryChanges(valuemap map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := dynamodb.New(sess)

	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		if len(value.(string)) != 0 {
			updateLanguage(svc, location, language, itemType, value.(string))
			logTransaction(svc, location, language, itemType, value.(string), "change")
		}
	}

}

func getInventoryLog(officeID string) []transactionLog {
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
	proj := expression.NamesList(expression.Name("Language"), expression.Name("TransactionID"),
		expression.Name("ItemType"), expression.Name("Location"), expression.Name("Transaction"),
		expression.Name("TransactionTime"), expression.Name("Count"))

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
		TableName:                 aws.String(transactionLogTableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	transactionList := make([]transactionLog, len(result.Items))

	for j, i := range result.Items {
		item := transactionLog{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		transactionList[j] = item
	}

	return transactionList
}

func AddInventoryChanges(valuemap map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	svc := dynamodb.New(sess)

	for key, value := range valuemap {
		s := strings.Split(key, "-")
		location := s[0]
		language := s[1]
		itemType := s[2]

		originalLocation := locationMaps[location].Languages[language+location]

		originalValue := originalLocation.Stats.ItemCounts[itemType]

		// in this case the office has a negative number of books.  It has placed
		// more books than it took in from the depot.  Maybe publishers are bringing
		// publications from home.
		if originalValue.(int) < 0 {
			originalValue = 0
		}

		if len(value.(string)) != 0 {
			newIntValue, err := strconv.Atoi(value.(string))
			if err != nil {
				fmt.Println("Query API call failed:")
				fmt.Println((err.Error()))
			}

			updateLanguage(svc, location, language, itemType, strconv.Itoa(originalValue.(int)+newIntValue))
			logTransaction(svc, location, language, itemType, strconv.Itoa(newIntValue),
				"add")
		}
	}
}

func SubtractInventoryFromDepot(valuemap map[string]interface{}) {
	getInventoryByOfficeID("Depot")
	originalLocation := locationMaps["Depot"]
	originalLanguages := originalLocation.Languages

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	svc := dynamodb.New(sess)

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
				updateLanguage(svc, "Depot", language, itemType, strconv.Itoa(originalValue.(int)-newIntValue))
				logTransaction(svc, "Depot", language, itemType, strconv.Itoa(newIntValue),
					"subtract")
			}
		}
	}
}
