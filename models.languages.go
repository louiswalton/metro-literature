package main

// import (
// 	"fmt"
// 	// 	//"os"
// 	// 	// "strconv"
// 	// 	// "strings"
// 	// 	// // "github.com/aws/aws-sdk-go/aws"
// 	// 	// "github.com/aws/aws-sdk-go/aws/session"
// 	// 	// "github.com/aws/aws-sdk-go/service/dynamodb"
// 	// 	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// 	// 	// "github.com/aws/aws-sdk-go/service/dynamodb/expression"
// )

type LocationStatsMap struct {
	ItemCounts map[string]interface{}
}

type language struct {
	Location       string `json:"location"`
	Language       string `json:"language"`
	Books          int    `json:"books"`
	BooksGoal      int    `json:"booksgoal"`
	Magazines      int    `json:"magazines"`
	MagazinesGoal  int    `json:"magazinesgoal"`
	Bibles         int    `json:"bibles"`
	BiblesGoal     int    `json:"biblesgoal"`
	DepotBooks     int
	DepotBibles    int
	DepotMagazines int
	Stats          *LocationStatsMap
}

type Location struct {
	Languages map[string]language
	Location  string
}

type lList []language

func (l lList) Len() int           { return len(l) }
func (l lList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l lList) Less(i, j int) bool { return l[i].Language > l[j].Language }

var locationMaps = make(map[string]Location)

func MergeDepotStats(lList []language) []language {
	getInventoryByOfficeID("Depot")
	depotInventory := locationMaps["Depot"]
	languageList := make([]language, len(lList))

	for i, lang := range lList {
		lang.DepotBooks = depotInventory.Languages[lang.Language+"Depot"].Books
		lang.DepotBibles = depotInventory.Languages[lang.Language+"Depot"].Bibles
		lang.DepotMagazines = depotInventory.Languages[lang.Language+"Depot"].Magazines
		languageList[i] = lang
	}
	return languageList
}
