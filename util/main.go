package util

import "log"

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
