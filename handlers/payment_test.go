package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/devenney/form3/common"
	"github.com/devenney/form3/payments"
	"github.com/devenney/form3/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kylelemons/godebug/pretty"
)

var mockPaymentList payments.PaymentList

func TestMain(m *testing.M) {
	// FIXME(devenney): This whole block should really be common code
	// to ensure we're testing the real server.
	common.InitConfig()

	var err error
	mockPaymentList, err = utils.InsertMockData("../payments/test_data/mock_data.json")
	if err != nil {
		os.Exit(1)
	}

	m.Run()
}

// addPayment adds a payment
func addPayment(ts *httptest.Server, payment *payments.Payment) (string, error) {
	url := ts.URL + "/addPayment"

	payload, err := json.Marshal(payment)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "encoding/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	if status := resp.StatusCode; status != http.StatusOK {
		return "", fmt.Errorf("Response has wrong status code: %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	id, err := uuid.Parse(string(body))
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

// getPayment gets a payment
func getPayment(ts *httptest.Server, id string) (p payments.Payment, err error) {
	url := ts.URL + "/getPayment/" + id

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if status := resp.StatusCode; status != http.StatusOK {
		err = fmt.Errorf("Response has wrong status code: %d, expected %d", resp.StatusCode, http.StatusOK)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var responsePayment payments.Payment
	err = json.Unmarshal(body, &responsePayment)
	if err != nil {
		return
	}

	p = responsePayment
	return
}

// updatePayment updates a payment
func updatePayment(ts *httptest.Server, payment *payments.Payment, id string) error {
	client := http.Client{}

	url := ts.URL + "/updatePayment/" + id

	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if status := resp.StatusCode; status != http.StatusOK {
		return fmt.Errorf("Response has wrong status code: %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	return nil
}

// deletePayment deletes a payment
func deletePayment(ts *httptest.Server, id string) error {
	url := ts.URL + "/deletePayment/" + id
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if status := resp.StatusCode; status != http.StatusOK {
		return fmt.Errorf("Response has wrong status code: %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	return nil
}

// TestGetPayment tests that each mock payment can be retrieved.
func TestGetPayment(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/getPayment/{id}", GetPaymentHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, mockPayment := range mockPaymentList.Data {
		responsePayment, err := getPayment(ts, mockPayment.ID)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(mockPayment, responsePayment) {
			diff := pretty.Compare(mockPayment, responsePayment)
			t.Fatalf("Mock payment and response payment were not equal: %s", diff)
		}
	}
}

// TestGetPayments tests listing a batch of payments, specifically
// the mock data set.
func TestGetPayments(t *testing.T) {
	req, err := http.NewRequest("GET", "/getPayments", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListPaymentsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var paymentList payments.PaymentList

	json.Unmarshal(rr.Body.Bytes(), &paymentList)

	for _, mockPayment := range mockPaymentList.Data {
		if !strings.Contains(rr.Body.String(), mockPayment.ID) {
			t.Fatalf("Response does not contain ID: %s.", mockPayment.ID)
		}
	}
}

// TestPaymentLifecycle tests the creation, reading and
// deletion of a payment.
func TestPaymentLifecycle(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/getPayment/{id}", GetPaymentHandler)
	r.HandleFunc("/addPayment", AddPaymentHandler)
	r.HandleFunc("/deletePayment/{id}", DeletePaymentHandler)
	r.HandleFunc("/updatePayment/{id}", UpdatePaymentHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payment := payments.Payment{
		PaymentType: "BACS",
	}

	id, err := addPayment(ts, &payment)
	if err != nil {
		t.Fatal(err)
	}

	updatedPayment, err := getPayment(ts, id)
	if err != nil {
		t.Fatal(err)
	}

	if updatedPayment.PaymentType != "BACS" {
		t.Fatalf("Response PaymentType was wrong: %s, expected: %s", updatedPayment.PaymentType, "BACS")
	}

	newPayment := payments.Payment{
		PaymentType: "FT",
	}

	err = updatePayment(ts, &newPayment, id)
}
