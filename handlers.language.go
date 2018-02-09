package main

import (
	"fmt"
	"net/http"
	// "reflect"
	// "strconv"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	officeID := "Depot"
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":    officeID,
			"location": officeID,
			"payload":  languages,
			"offices":  offices,
		},
	)

}

func showInventoryByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":    officeID,
			"location": officeID,
			"payload":  languages,
			"offices":  offices,
		},
	)

}

func editInventoryByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()

	c.HTML(
		http.StatusOK,
		"editOfficeInventory.html",
		gin.H{
			"title":    officeID,
			"location": officeID,
			"payload":  languages,
			"offices":  offices,
		},
	)
}

func addInventoryByOfficeId(c *gin.Context) {
	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	depotInventory := getInventoryByOfficeID("Depot")
	offices := getAllOffices()

	c.HTML(
		http.StatusOK,
		"addOfficeInventory.html",
		gin.H{
			"title":    officeID,
			"location": officeID,
			"payload":  languages,
			"offices":  offices,
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

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":    officeID,
			"location": officeID,
			"payload":  languages,
			"offices":  offices,
		},
	)
}

func addInventoryToOffice(c *gin.Context) {
	c.Request.ParseForm()
	valueMap := make(map[string]interface{})
	for key, value := range c.Request.PostForm {
		valueMap[key] = value[0]
	}

	AddInventoryChanges(valueMap)

	officeID := c.Param("office_id")
	languages := getInventoryByOfficeID(officeID)
	offices := getAllOffices()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"officeInventory.html",
		// Pass the data that the page uses
		gin.H{
			"title":    officeID,
			"location": officeID,
			"payload":  languages,
			"offices":  offices,
		},
	)
}
