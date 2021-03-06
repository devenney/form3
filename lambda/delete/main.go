package main

import (
	"log"

	"github.com/apex/gateway"
	"github.com/devenney/form3/common"
	"github.com/devenney/form3/handlers"
	"github.com/gorilla/mux"
)

func main() {
	common.InitConfig()

	r := mux.NewRouter()

	r.HandleFunc("/deletePayment/{id}", handlers.DeletePaymentHandler)

	log.Fatal(gateway.ListenAndServe(":3000", r))
}
