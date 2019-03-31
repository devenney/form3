package database

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/viper"
)

const (
	paymentsTable = "payments"
)

var instance *Service
var mu sync.Mutex

// GetInstance returns the Service singleton,
// initialising it if it does not exists.
func GetInstance() *Service {
	mu.Lock() // FIXME(devenney): Unnecessary locking.
	defer mu.Unlock()

	if instance == nil {
		instance = &Service{
			endpoint: viper.GetString("db_endpoint"),
		}
		err := instance.InitDB()
		if err != nil {
			panic(err)
		}
	}
	return instance
}

// Service represents an established DynamoDB service
type Service struct {
	endpoint string
	Svc      *dynamodb.DynamoDB
}

// InitDB initialises a Service, setting up the
// necessary connections and creating tables if they
// do not exist already.
func (db *Service) InitDB() error {
	config := &aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String(db.endpoint),
	}

	if viper.Get("env") != "PROD" {
		config.Credentials = credentials.NewStaticCredentials("AKID", "SECRET", "TOKEN")
	}

	sess := session.Must(session.NewSession(config))

	db.Svc = dynamodb.New(sess)

	var tableMap = map[string]func() error{
		paymentsTable: db.createPaymentsTable,
	}

	for tableName, createFunc := range tableMap {
		exists, err := db.tableExists(tableName)
		if err != nil {
			return err
		}

		if !exists {
			err = createFunc()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
