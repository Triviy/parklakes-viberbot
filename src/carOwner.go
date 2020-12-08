package main

import (
	"regexp"
	"strings"
)

// CarOwner represents basic car owner data
type CarOwner struct {
	ID        string   `bson:"_id,omitempty"`
	Created   string   `bson:"created"`
	CarNumber string   `bson:"carNumber"`
	Owner     string   `bson:"owner"`
	Phones    []string `bson:"phones"`
}

// CarOwnerProp represents key/value property
type CarOwnerProp struct {
	ID    string `bson:"_id,omitempty"`
	Value string `bson:"value"`
}

var replacer = strings.NewReplacer(
	"А", "A",
	"В", "B",
	"С", "C",
	"Е", "E",
	"Н", "H",
	"І", "I",
	"К", "K",
	"М", "M",
	"О", "O",
	"Р", "P",
	"Т", "T",
	"И", "U",
	"Х", "X",
	"У", "Y",
	" ", "",
	"-", "",
)
var standardNumberRegexp = regexp.MustCompile(`[A-Z]{2}[0-9]{4}[A-Z]{2}|[0-9]{5}[A-Z,А-Я,І]{2}|[A-Z]{3}[0-9]{3}`)

// CreateCarOwnerFromRecord creates *CarOwner struct from a record
func CreateCarOwnerFromRecord(record []interface{}) *CarOwner {
	carNumber := strings.TrimSpace(record[1].(string))
	if carNumber == "" {
		return nil
	}
	carOwner := CarOwner{
		ID:        NormalizeCarNumber(carNumber),
		CarNumber: carNumber,
		Created:   record[0].(string),
		Owner:     record[2].(string),
	}
	firstPhone := strings.TrimSpace(record[3].(string))
	if firstPhone != "" {
		carOwner.Phones = append(carOwner.Phones, firstPhone)
	}
	if len(record) == 5 {
		secondPhone := strings.TrimSpace(record[4].(string))
		if secondPhone != "" {
			carOwner.Phones = append(carOwner.Phones, secondPhone)
		}
	}
	return &carOwner
}

// NormalizeCarNumber returns searchable car number string
func NormalizeCarNumber(cn string) string {
	cn = strings.ToUpper(cn)
	cn = replacer.Replace(cn)
	matchedCn := standardNumberRegexp.FindString(cn)
	if matchedCn != "" {
		return matchedCn
	}
	return cn
}
