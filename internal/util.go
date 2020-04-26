package internal

import "log"

func FailOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LogError(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}
