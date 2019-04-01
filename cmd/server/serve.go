// Package main runs the monolithic API server, primarily
// intended for easy local development.
package main

import (
	"github.com/devenney/form3/server"
)

func main() {
	server.StartServer()
}
