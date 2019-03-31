package database

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/devenney/form3/common"
)

// TestTableExists tests tableExists in positive
// and negative cases
func TestTableExists(t *testing.T) {
	input := &dynamodb.CreateTableInput{
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
		TableName: aws.String("test_table"),
	}

	db := GetInstance()
	db.Svc.CreateTable(input)

	// Test that our table has created, and our code sees this
	exists, err := db.tableExists("test_table")
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatalf("test_table was not created, no error was produced")
	}

	// Test that a non-existence table is not reported as existing
	exists, err = db.tableExists("not_test_table")
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatalf("not_test_table unexpectedly exists")
	}
}

// TestTableCreation tests all of our table creation functions
func TestTableCreation(t *testing.T) {
	common.InitConfig()

	db := GetInstance()

	// A map of table name to creation function,
	// as in database.go
	var tableMap = map[string]func() error{
		paymentsTable: db.createPaymentsTable,
	}

	// Loop through our creation funcs, execute and test
	for _, createFunc := range tableMap {
		err := createFunc()
		if err != nil {
			t.Fatal(err)
		}
	}
}
