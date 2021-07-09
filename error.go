package main

import (
	"fmt"
	"log"
)

func CheckIfError(err error) {
	if err == nil {
		return
	}

	log.Fatal(fmt.Sprintf("error: %s", err))
}
