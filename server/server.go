// Package server runs the monolithic API server, primarily
// intended for easy local development.
package server

import (
	"log"
	"net/http"

	"github.com/devenney/form3/common"
	"github.com/devenney/form3/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func commonMiddleware(next http.Handler) http.Handler {
	// TODO(devenney): Auth.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// StartServer starts a monolithic server
func StartServer() {
	common.InitConfig()

	r := mux.NewRouter()
	r.Use(commonMiddleware)

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
