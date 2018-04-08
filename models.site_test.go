package main

import "testing"

// Test the function that gets all of the languages
func TestCreateLocationReport(t *testing.T) {
	m := make(map[string]int)
	m["English"] = 9
	m["Chinese"] = 3
	m["Day"] = 7
	m["Month"] = 4
	m["Year"] = 2018

	s := site{
		Location: "test Location",
		Office:   "test Office",
	}

	s.CreateReport(m)
}

func TestGetLocationReport(t *testing.T) {
	m := make(map[string]int)
	m["English"] = 9
	m["Chinese"] = 3

	s := site{
		Location: "getReportLocation",
		Office:   "getReportOffice",
	}

	s.CreateReport(m)

	reports := s.GetLocationReports("getReportLocation")

	if len(reports) != 1 {
		t.Fail()
	}

	var report = reports[0]
	if report.Report["English"] != 9 {
		t.Fail()
	}
}