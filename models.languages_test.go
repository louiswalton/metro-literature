package main

import "testing"

// Test the function that gets all of the languages
func TestGetAllLanguages(t *testing.T) {
	lList := getAllLanguages()

	if len(lList) != len(languageList) {
		t.Fail()
	}

	for i, v := range lList {
		if v.Books != languageList[i].Books ||
			v.Name != languageList[i].Name ||
			v.Magazines != languageList[i].Magazines ||
			v.Bibles != languageList[i].Bibles {
			t.Fail()
			break
		}
	}
}
