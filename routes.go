// routes.go
package main

func initializeRoutes() {
	//Handle the index route
	router.GET("/", showIndexPage)
	router.GET("/:office_id", showInventoryByOfficeId)
	router.GET("/:office_id/editInventory", editInventoryByOfficeId)
	router.GET("/:office_id/addInventory", addInventoryByOfficeId)
	router.POST("/:office_id/editInventory", saveInventoryChanges)
	router.POST("/:office_id/addInventory", addInventoryToOffice)
	// router.GET("/:office_id/reservation", getStockReservation)
	// router.POST("/:office_id/reservation", postStockReservation)
	// router.PUT("/:office_id/reservation", putStockReservation)
	// router.GET("/:office_id/order_history", getOrderHistory)
	// router.POST("/:location/report", postLocationReport)

}
