package models

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
	"\t", "",
	"\n", "",
)
var standardNumberRegexp = regexp.MustCompile(`[A-Z]{2}[0-9]{4}[A-Z]{2}|[0-9]{5}[A-Z,А-Я,І]{2}|[A-Z]{3}[0-9]{3}`)

// NormalizeCarNumber returns searchable car number string
func NormalizeCarNumber(cn string) string {
	cn = strings.ToUpper(cn)
	cn = replacer.Replace(cn)
	if matchedCn := standardNumberRegexp.FindString(cn); matchedCn != "" {
		return matchedCn
	}
	return cn
}
