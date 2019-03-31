package payments

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/devenney/form3/common"
	jd "github.com/josephburnett/jd/lib"
)

const (
	mockDataFile = "test_data/mock_data.json"
)

// TestUnmarshal tests the ability to unmarshal mock data into a struct
// and marshal it back to JSON without data loss.
func TestUnmarshal(t *testing.T) {
	// Read mock data
	file, err := ioutil.ReadFile(mockDataFile)
	if err != nil {
		t.Fatalf("Failed to read mock data file %s: %v", mockDataFile, err)
	}

	// Unmarshal our mock data to a payment list
	var paymentList PaymentList
	err = json.Unmarshal(file, &paymentList)
	if err != nil {
		t.Logf("File contents: %s", file)
		t.Fatalf("Failed to unmarshal mock data: %v", err)
	}

	// Re-marshal the payment list to test data loss
	marshalled, err := json.Marshal(paymentList)
	if err != nil {
		t.Fatalf("Failed to re-marshal mock data: %v", err)
	}

	// Read mock data into JSON differ
	a, err := jd.ReadJsonString(string(file))
	if err != nil {
		t.Fatalf("Failed to read JSON string from file: %s", file)
	}

	// Read re-marshalled data into JSON differ
	b, err := jd.ReadJsonString(string(marshalled))
	if err != nil {
		t.Fatalf("Failed to read marshalled JSON string: %s", marshalled)
	}

	// Diff the mock data and the re-marshalled data
	diff := a.Diff(b)
	if len(diff) != 0 {
		t.Fatalf("Original and marshalled JSON have diverged:: %s", diff.Render())
	}
}

// TestInsertAndDelete inserts then deletes a Payment
func TestInsertAndDelete(t *testing.T) {
	common.InitConfig()

	payment := Payment{
		PaymentType: "BACS",
	}

	// Insert
	err := payment.Upsert()
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the inserted Payment
	paymentGet, err := Get(payment.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the ID is as expected
	if paymentGet.ID != payment.ID {
		t.Fatalf("ID changed between write and read: %s, expected: %s", paymentGet.ID, payment.ID)
	}

	// Delete
	err = payment.Delete()
	if err != nil {
		t.Fatal(err)
	}
}
