// Package main runs the monolithic API server, primarily
// intended for easy local development.
package main

import (
	"log"
	"net/http"

	"github.com/devenney/form3/common"
	"github.com/devenney/form3/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	common.InitConfig()

	r := mux.NewRouter()

	// Add our endpoints, restricted by method
	// NOTE: Make sure this tracks the API Gateway methods in serverless.yml
	r.HandleFunc("/getPayments", handlers.ListPaymentsHandler).Methods("GET")
	r.HandleFunc("/getPayment/{id}", handlers.GetPaymentHandler).Methods("GET")
	r.HandleFunc("/addPayment", handlers.AddPaymentHandler).Methods("POST")
	r.HandleFunc("/updatePayment/{id}", handlers.UpdatePaymentHandler).Methods("PUT")
	r.HandleFunc("/deletePayment/{id}", handlers.DeletePaymentHandler).Methods("DELETE")

	srv := &http.Server{
		Handler: r,
		Addr:    viper.GetString("BIND_ADDRESS"),
	}

	// Serve
	log.Printf("Booting server, listening on %s", viper.GetString("BIND_ADDRESS"))
	log.Fatal(srv.ListenAndServe())
}
