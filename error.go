package main

import (
	"fmt"
	"log"
)

func checkIfError(err error) {
	if err == nil {
		return
	}

	log.Fatal(fmt.Sprintf("error: %s", err))
}
