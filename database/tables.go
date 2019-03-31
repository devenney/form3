package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// tableExists tests if a table already exists. We default to false
// if an error other than ErrCodeResourceNotFoundException is
// encountered, so the caller should always assume the result
// is unknown if there is an error returned.
func (db *Service) tableExists(tableName string) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	_, err := db.Svc.DescribeTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				return false, err
			}
		} else {
			return false, err
		}
	}
	return true, nil
}

// createPaymentsTable creates a table to store payments
func (db *Service) createPaymentsTable() error {
	createPaymentsTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(paymentsTable),
	}

	// TODO(devenney): Error?
	db.Svc.CreateTable(createPaymentsTableInput)
	return nil
}
