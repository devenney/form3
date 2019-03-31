package common

import (
	"log"

	"github.com/spf13/viper"
)

const (
	envPrefix = "PAYMENTS"

	// APIUrlKey is the API URL config key
	APIUrlKey     = "API_URL"
	apiURLDefault = "localhost"

	// BindAddressKey is the HTTP bind config key
	BindAddressKey     = "BIND_ADDRESS"
	bindAddressDefault = "127.0.0.1:8000"

	// DBEndpointKey is the Dynamo Endpoint config key
	DBEndpointKey     = "DB_ENDPOINT"
	dbEndpointDefault = "http://localhost:8080"

	// EnvKey is the deployment environemtn config key
	EnvKey     = "ENV"
	envDefault = "LOCAL"
)

// InitConfig initialises the application config, setting
// development defaults which can be overridded by environment
// variables.
func InitConfig() {
	viper.SetEnvPrefix(envPrefix)

	// API URL
	viper.SetDefault(APIUrlKey, apiURLDefault)
	viper.BindEnv(APIUrlKey)

	// Bind Address
	viper.SetDefault(BindAddressKey, bindAddressDefault)
	viper.BindEnv(BindAddressKey)

	// Database Endpoints
	viper.SetDefault(DBEndpointKey, dbEndpointDefault)
	viper.BindEnv(DBEndpointKey)

	// Environment
	viper.SetDefault(EnvKey, envDefault)
	viper.BindEnv(EnvKey)

	// Development logging
	if viper.GetString(EnvKey) != "PROD" {
		logConfig()
	}
}

// logConfig logs the interesting configuration values for debug
func logConfig() {
	log.Printf("Viper | %s is %s", APIUrlKey, viper.GetString(APIUrlKey))
	log.Printf("Viper | %s is %s", BindAddressKey, viper.GetString(BindAddressKey))
	log.Printf("Viper | %s is %s", DBEndpointKey, viper.GetString(DBEndpointKey))
	log.Printf("Viper | %s is %s", EnvKey, viper.GetString(EnvKey))
}
