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

func LogInfo(err error) {
	if err != nil {
		log.Printf("INFO: %s", err)
	}
}

func LogDebug(err error) {
	if err != nil {
		log.Printf("DEBUG: %s", err)
	}
}
