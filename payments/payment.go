// Package payments holds the modelss and
// business logic related to payments.
package payments

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/devenney/form3/database"
	"github.com/go-validator/validator"
	"github.com/google/uuid"
)

// PaymentList holds a list of Payments
type PaymentList struct {
	Data  []Payment    `json:"data"`
	Links PaymentLinks `json:"links" validate:"nonzero"`
}

// Validate validates each payment in a PaymentList
func (pl *PaymentList) Validate() error {
	for _, p := range pl.Data {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Payment represents an individual payment
type Payment struct {
	Attributes     PaymentAttributes `json:"attributes" validate:"nonzero"`
	ID             string            `json:"id"`
	OrganisationID string            `json:"organisation_id" validate:"nonzero"`
	PaymentType    string            `json:"type" validate:"nonzero"`
	Version        int               `json:"version" validate:"min=0"`
}

// Validate validates a payment based on struct tags
func (p *Payment) Validate() error {
	return validator.Validate(p)
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
	Amount               string                    `json:"amount,float" validate:"nonzero"`
	BeneficiaryParty     PaymentParty              `json:"beneficiary_party" validate:"nonzero"`
	ChargesInformation   PaymentChargesInformation `json:"charges_information" validate:"nonzero"`
	Currency             string                    `json:"currency" validate:"nonzero"`
	DebtorParty          PaymentParty              `json:"debtor_party" validate:"nonzero"`
	EndToEndReference    string                    `json:"end_to_end_reference" validate:"nonzero"`
	Fx                   PaymentForeignExchange    `json:"fx"`
	NumericReference     int                       `json:"numeric_reference,string" validate:"min=0"`
	PaymentID            string                    `json:"payment_id" validate:"nonzero"`
	PaymentPurpose       string                    `json:"payment_purpose" validate:"nonzero"`
	PaymentScheme        string                    `json:"payment_scheme" validate:"nonzero"`
	PaymentType          string                    `json:"payment_type" validate:"nonzero"`
	ProcessingDate       string                    `json:"processing_date" validate:"nonzero"`
	Reference            string                    `json:"reference" validate:"nonzero"`
	SchemePaymentSubType string                    `json:"scheme_payment_sub_type" validate:"nonzero"`
	SchemePaymentType    string                    `json:"scheme_payment_type" validate:"nonzero"`
	SponsorParty         PaymentParty              `json:"sponsor_party" validate:"nonzero"`
}

// PaymentParty represents a party (e.g. Beneficiary, Debtor)
// involved in a payment.
type PaymentParty struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty" validate:"nonzero"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	AccountType       *int   `json:"account_type,omitempty"`
	Address           string `json:"address,omitempty"`
	BankID            int    `json:"bank_id,string,omitempty" validate:"min=0"`
	BankIDCode        string `json:"bank_id_code,omitempty" validate:"nonzero"`
	Name              string `json:"name,omitempty" pact:"name"`
}

// PaymentChargesInformation represents a set of charges
// associated with a Payment.
type PaymentChargesInformation struct {
	BearerCode              string           `json:"bearer_code" validate:"nonzero"`
	ReceiverChargesAmount   string           `json:"receiver_charges_amount,float" validate:"nonzero"`
	ReceiverChargesCurrency string           `json:"receiver_charges_currency" validate:"nonzero"`
	SenderCharges           []PaymentCharges `json:"sender_charges" validate:"nonzero"`
}

// PaymentCharges represents the actual charge associated
// with a payment, including its currency.
type PaymentCharges struct {
	Amount   string `json:"amount,float" validate:"nonzero"`
	Currency string `json:"currency" validate:"nonzero"`
}

// PaymentForeignExchange represents FX details when
// a payment has been converted between currencies.
type PaymentForeignExchange struct {
	ContractReference string `json:"contract_reference" validate:"nonzero"`
	ExchangeRate      string `json:"exchange_rate,float" validate:"nonzero"`
	OriginalAmount    string `json:"original_amount,float" validate:"nonzero"`
	OriginalCurrency  string `json:"original_currency" validate:"nonzero"`
}
