package main

import "testing"

// Test the function that gets all of the languages
func TestGetStockListNotEmpty(t *testing.T) {
	var sList = getStocklist()

	if len(sList.Stocklist) == 0 {
		t.Fail()
	}

}

func TestStockListSanity(t *testing.T) {
	var sList = getStocklist()

	_, ok := sList.Stocklist["English"]

	if !ok {
		t.Fail()
	}
}
