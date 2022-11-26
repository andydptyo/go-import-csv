package main

import (
	"log"

	"github.com/andydptyo/go-import-csv/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("error execute command %v", err)
	}
}
