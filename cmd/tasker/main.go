// Package main provides the entry point for the Tasker CLI application.
package main

import (
	"log"

	"tasker.jsas.dev/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
