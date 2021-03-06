package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/devenney/form3/database"
	"github.com/devenney/form3/payments"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

const (
	tableName = "payments"
)

// ListPaymentsHandler handles /getPayments
//
// TODO(devenney): Filter the payments based on some key. Presumably
//                 we would never want to just return all payments.
func ListPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Scan the entire table
	result, err := db.Svc.Scan(input)
	if err != nil {
		fmt.Println(w, err.Error())
	}

	// Build our list of payments
	var paymentList payments.PaymentList
	for _, result := range result.Items {
		var payment payments.Payment

		err = dynamodbattribute.UnmarshalMap(result, &payment)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("Unable to unmarshal Dynamo item: %v", err.Error())
		}

		paymentList.Data = append(paymentList.Data, payment)
	}

	// Insert our API links
	paymentList.Links = payments.PaymentLinks{
		Self: viper.GetString("API_URL"),
	}

	output, err := json.Marshal(paymentList)
	if err != nil {
		log.Printf("Unable to marshal payment list: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(output))
}

// GetPaymentHandler handles /getPayment/{id}
func GetPaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	payment, err := payments.Get(vars["id"])
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check for non-existence
	if payment.ID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Marshal our object
	output, err := json.Marshal(payment)
	if err != nil {
		log.Printf("Unable to marshal payment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(output))
}

// AddPaymentHandler handles /addPayment
//
// TODO(devenney): Input validation. We will currently accept an
//                 empty payload and return a UUID.
// TODO(devenney): We should be inserting here, not upserting.
//                 Currently we will overwrite objects.
func AddPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var payment payments.Payment

	// Read HTTP body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Unmarshal body to a Payment
	err = json.Unmarshal(payload, &payment)
	if err != nil {
		log.Printf("Unable to Unmarshal payload: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Ensure the user has not self-set an invalid ID
	if payment.ID != "" {
		_, err = uuid.Parse(payment.ID)
		if err != nil {
			log.Printf("Payment ID was not a valid UUID: %s, %v", payment.ID, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	// Validate
	if err := payment.Validate(); err != nil {
		log.Printf("Failed to validate Payment: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Insert
	err = payment.Upsert()
	if err != nil {
		log.Printf("Error upserting Payment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Return the ID
	fmt.Fprintf(w, payment.ID)
}

// UpdatePaymentHandler handles /updatePayment
//
// TODO(devenney): Input validation. We will currently accept an
//                 empty payload and write it to the database.
func UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	payment, err := payments.Get(vars["id"])
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check for non-existence
	if payment.ID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var updatedPayment payments.Payment

	// Read HTTP body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Unmarshal body to a Payment
	err = json.Unmarshal(payload, &updatedPayment)
	if err != nil {
		log.Printf("Unable to Unmarshal payload: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// If the user didn't specify the ID in the body,
	// we insert it here to avoid autogenerating a new ID.
	if updatedPayment.ID == "" {
		updatedPayment.ID = payment.ID
	}

	// If the user did specify the ID in the body, make sure
	// it matches the path variable.
	if payment.ID != updatedPayment.ID {
		log.Printf("Provided payment ID (%s) did not match payment to be updated (%s)", payment.ID, updatedPayment.ID)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Validate
	if err := payment.Validate(); err != nil {
		log.Printf("Failed to validate Payment: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Update
	err = updatedPayment.Upsert()
	if err != nil {
		log.Printf("Error upserting Payment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Return the ID
	fmt.Fprintf(w, payment.ID)
}

// DeletePaymentHandler handles /deletePayment/{id}
func DeletePaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	payment, err := payments.Get(vars["id"])
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check for non-existence
	if payment.ID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Delete
	err = payment.Delete()
	if err != nil {
		log.Printf("Failed to delete payment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
