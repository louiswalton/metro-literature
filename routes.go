// routes.go
package main

func initializeRoutes() {
	//Handle the index route
	router.GET("/", showIndexPage)
	router.GET("/office/:office_id", showInventoryByOfficeId)
	router.GET("/office/:office_id/editInventory", editInventoryByOfficeId)
	router.GET("/office/:office_id/addInventory", addInventoryByOfficeId)
	router.POST("/office/:office_id/editInventory", saveInventoryChanges)
	router.POST("/office/:office_id/addInventory", addInventoryToOffice)
	router.GET("/office/:office_id/inventoryLog", showInventoryLogByOfficeId)
	// router.GET("/office/:office_id/reservation", getStockReservation)
	// router.POST("/office/:office_id/reservation", postStockReservation)
	// router.PUT("/office/:office_id/reservation", putStockReservation)
	// router.GET("/office/:office_id/order_history", getOrderHistory)
	// router.POST("/Location/:location", postLocationReport)
	router.GET("/location/:office_id/:location", showLocationReport)
	router.GET("/location/:office_id/:location/createReport", createLocationReport)
	router.POST("/location/:office_id/:location/saveReport", saveLocationReport)
}
