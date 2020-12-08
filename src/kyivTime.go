package main

import (
	"log"
	"time"
)

const kyiv = "Europe/Kiev"
const kyivFormat = "02.01.2006 15:04:05"

// ToKyivTime converts string value to Kyiv time
func ToKyivTime(s string) time.Time {
	loc, err := time.LoadLocation(kyiv)
	if err != nil {
		log.Fatalf("Error while getting location: %v", err)
	}
	t, err := time.ParseInLocation(kyivFormat, s, loc)
	if err != nil {
		log.Fatalf("Error while parsing time (%s): %v", s, err)
	}
	return t
}

// GetKyivTime returns current time.Time in Kyiv
func GetKyivTime() time.Time {
	loc, err := time.LoadLocation(kyiv)
	if err != nil {
		log.Fatalf("Error while getting location: %v\n", err)
	}
	return time.Now().In(loc)
}
