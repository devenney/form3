package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/devenney/form3/payments"
)

// InsertMockData populates the database with mock payments from a file.
func InsertMockData(path string) (paymentList payments.PaymentList, err error) {
	// Read mock data file
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read mock data file %s: %v", path, err)
	}

	// Unmarshal the list of payments
	err = json.Unmarshal(f, &paymentList)
	if err != nil {
		log.Fatalf("Failed to unmarshal mock data: %v", err)
	}

	// Upsert each payment
	for _, payment := range paymentList.Data {
		err = payment.Upsert()
		if err != nil {
			return
		}
	}

	return
}
