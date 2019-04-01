// Package payments holds the modelss and
// business logic related to payments.
package payments

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/devenney/form3/database"
	"github.com/google/uuid"
)

// PaymentList holds a list of Payments
type PaymentList struct {
	Data  []Payment    `json:"data"`
	Links PaymentLinks `json:"links,omitempty"`
}

// Payment represents an individual payment
type Payment struct {
	Attributes     PaymentAttributes `json:"attributes"`
	ID             string            `json:"id"`
	OrganisationID string            `json:"organisation_id"`
	PaymentType    string            `json:"type"`
	Version        int               `json:"version"`
}

// Get queries DynamoDB for a payment with the given ID
func Get(id string) (*Payment, error) {
	db := database.GetInstance()

	input := &dynamodb.GetItemInput{
		TableName: aws.String("payments"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := db.Svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	var payment Payment
	err = dynamodbattribute.UnmarshalMap(result.Item, &payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// Upsert updates a Payment into the database, creating it
// in the process if it does not exist.
func (p *Payment) Upsert() error {
	db := database.GetInstance()

	if p.ID == "" {
		p.ID = uuid.Must(uuid.NewRandom()).String()
	}

	av, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("payments"),
	}

	_, err = db.Svc.PutItem(input)
	return err
}

// Delete removes a payment from the database.
func (p *Payment) Delete() error {
	db := database.GetInstance()

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("payments"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(p.ID),
			},
		},
	}

	_, err := db.Svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

// PaymentLinks holds the relevant links to return
// to clients after API calls.
type PaymentLinks struct {
	Self string `json:"self,omitempty"`
}

// PaymentAttributes holds the attributes of an individual
// payment. Floats are stored as strings to prevent data loss.
type PaymentAttributes struct {
	Amount               string                    `json:"amount,float"`
	BeneficiaryParty     PaymentParty              `json:"beneficiary_party"`
	ChargesInformation   PaymentChargesInformation `json:"charges_information"`
	Currency             string                    `json:"currency"`
	DebtorParty          PaymentParty              `json:"debtor_party"`
	EndToEndReference    string                    `json:"end_to_end_reference"`
	Fx                   PaymentForeignExchange    `json:"fx"`
	NumericReference     int                       `json:"numeric_reference,string"`
	PaymentID            string                    `json:"payment_id"`
	PaymentPurpose       string                    `json:"payment_purpose"`
	PaymentScheme        string                    `json:"payment_scheme"`
	PaymentType          string                    `json:"payment_type"`
	ProcessingDate       string                    `json:"processing_date"`
	Reference            string                    `json:"reference"`
	SchemePaymentSubType string                    `json:"scheme_payment_sub_type"`
	SchemePaymentType    string                    `json:"scheme_payment_type"`
	SponsorParty         PaymentParty              `json:"sponsor_party"`
}

// PaymentParty represents a party (e.g. Beneficiary, Debtor)
// involved in a payment.
type PaymentParty struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	AccountType       *int   `json:"account_type,omitempty"`
	Address           string `json:"address,omitempty"`
	BankID            int    `json:"bank_id,string,omitempty"`
	BankIDCode        string `json:"bank_id_code,omitempty"`
	Name              string `json:"name,omitempty" pact:"name"`
}

// PaymentChargesInformation represents a set of charges
// associated with a Payment.
type PaymentChargesInformation struct {
	BearerCode              string           `json:"bearer_code"`
	ReceiverChargesAmount   string           `json:"receiver_charges_amount,float"`
	ReceiverChargesCurrency string           `json:"receiver_charges_currency"`
	SenderCharges           []PaymentCharges `json:"sender_charges"`
}

// PaymentCharges represents the actual charge associated
// with a payment, including its currency.
type PaymentCharges struct {
	Amount   string `json:"amount,float"`
	Currency string `json:"currency"`
}

// PaymentForeignExchange represents FX details when
// a payment has been converted between currencies.
type PaymentForeignExchange struct {
	ContractReference string `json:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate,float"`
	OriginalAmount    string `json:"original_amount,float"`
	OriginalCurrency  string `json:"original_currency"`
}
