package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

// init logs
func logInit() {
	os.MkdirAll("logs", 0755)
	ioutil.WriteFile("logs/warnings.log", []byte(""), 0755)
	ioutil.WriteFile("logs/errors.log", []byte(""), 0755)
}

// log warning, like product not exist, or no price
func logWarning(warning string) {
	f, err := os.OpenFile("logs/warnings.log", os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Error, can't open 'logs/warnings.log' file")
	}
	defer f.Close()

	msg := time.Now().Format("02/01/2006 15:04:05") + " " + warning + "\n\n"

	if _, err = f.WriteString(msg); err != nil {
		log.Fatal("Error, can't write to 'logs/warnings.log' file")
	}
}

// TODO
// log errors, like no response from server
func logError() {

}
