package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	officeID := "Depot"
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	fmt.Println("finished getSiteByOfficeID")

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":     officeID,
			"office":    officeID,
			"location":  officeID,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)

}

func showInventoryByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":     officeID,
			"office":    officeID,
			"location":  officeID,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)

}

func showInventoryLogByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	logs := getInventoryLog(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"showOfficeLogs.html",
		// Pass the data that the page uses
		gin.H{
			"title":     officeID,
			"office":    officeID,
			"location":  officeID,
			"locations": locations,
			"payload":   logs,
			"offices":   offices,
		},
	)
}

func showLocationReport(c *gin.Context) {
	officeID := c.Param("office_id")
	location := c.Param("location")
	logs := GetLocationReports(location)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)
	// 	fmt.Println(logs)
	fmt.Println("Get Location report:")
	fmt.Println(logs)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"locationLog.html",
		// Pass the data that the page uses
		gin.H{
			"office":    officeID,
			"title":     location,
			"location":  location,
			"locations": locations,
			"payload":   logs,
			"offices":   offices,
		},
	)
}

func createLocationReport(c *gin.Context) {
	officeID := c.Param("office_id")
	location := c.Param("location")
	languages := getInventoryByOfficeID(officeID)
	logs := GetLocationReports(location)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)
	// 	fmt.Println(logs)
	fmt.Println("Get Location report:")
	fmt.Println(logs)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"createLocationReport.html",
		// Pass the data that the page uses
		gin.H{
			"office":    officeID,
			"title":     location,
			"location":  location,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)
}

func editInventoryByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	c.HTML(
		http.StatusOK,
		"editOfficeInventory.html",
		gin.H{
			"title":     officeID,
			"office":    officeID,
			"location":  officeID,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)
}

func addInventoryByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	if officeID != "Depot" {
		languages = MergeDepotStats(languages)
	}

	c.HTML(
		http.StatusOK,
		"addOfficeInventory.html",
		gin.H{
			"title":     officeID,
			"office":    officeID,
			"location":  officeID,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)
}

func saveInventoryChanges(c *gin.Context) {
	c.Request.ParseForm()
	valueMap := make(map[string]interface{})
	for key, value := range c.Request.PostForm {
		valueMap[key] = value[0]
	}

	MakeInventoryChanges(valueMap)

	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":     officeID,
			"location":  officeID,
			"office":    officeID,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)
}

func saveLocationReport(c *gin.Context) {
	c.Request.ParseForm()
	valueMap := make(map[string]string)
	for key, value := range c.Request.PostForm {
		if value[0] != "" {
			valueMap[key] = fmt.Sprintf("%v", value[0])
		}
	}

	office := c.Param("office_id")
	location := c.Param("location")
	date := valueMap["reportDate"]
	delete(valueMap, "reportDate")
	logs := GetLocationReports(location)
	offices := getAllOffices()
	locations := getSiteByOfficeID(office)

	CreateReport(valueMap, location, office, date)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"locationLog.html",
		// Pass the data that the page uses
		gin.H{
			"office":    office,
			"title":     location,
			"location":  location,
			"locations": locations,
			"payload":   logs,
			"offices":   offices,
		},
	)
}

func saveEditLocationReport(c *gin.Context) {
	c.Request.ParseForm()
	valueMap := make(map[string]string)

	office := c.Param("office_id")
	location := c.Param("location")
	logs := GetLocationReports(location)
	offices := getAllOffices()
	locations := getSiteByOfficeID(office)
	var date = ""

	CreateReport(valueMap, location, office, date)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"locationLog.html",
		// Pass the data that the page uses
		gin.H{
			"office":    office,
			"title":     location,
			"location":  location,
			"locations": locations,
			"payload":   logs,
			"offices":   offices,
		},
	)
}

func editLocationReport(c *gin.Context) {
	c.Request.ParseForm()

	office := c.Param("office_id")
	location := c.Param("location")
	reportID := c.Param("report_id")
	offices := getAllOffices()
	locations := getSiteByOfficeID(office)

	report := GetReport(reportID)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"ReportView.html",
		// Pass the data that the page uses
		gin.H{
			"office":    office,
			"title":     location,
			"location":  location,
			"locations": locations,
			"report":    report,
			"offices":   offices,
		},
	)
}

func addInventoryToOffice(c *gin.Context) {
	c.Request.ParseForm()
	valueMap := make(map[string]interface{})
	for key, value := range c.Request.PostForm {
		valueMap[key] = value[0]
	}

	SubtractInventoryFromDepot(valueMap)
	AddInventoryChanges(valueMap)

	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()
	locations := getSiteByOfficeID(officeID)

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":     officeID,
			"office":    officeID,
			"location":  officeID,
			"locations": locations,
			"payload":   languages,
			"offices":   offices,
		},
	)
}
