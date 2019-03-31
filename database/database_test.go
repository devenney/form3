package database

import (
	"testing"
)

func TestGetInstance(t *testing.T) {
	db := GetInstance()

	if db == (&Service{}) {
		t.Fatal("Database was not initialised")
	}
}
