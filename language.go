package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	languages := getAllLanguages()
	offices := getAllOffices()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Home Page",
		"payload": languages,
		"offices": offices,
	})
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
