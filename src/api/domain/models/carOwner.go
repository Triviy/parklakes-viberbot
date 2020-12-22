package models

import (
	"fmt"
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
var standardNumberRegexp = regexp.MustCompile(`[A-Z]{2}[0-9]{4}[A-Z]{2}|[0-9]{5}[A-Z,А-Я,І]{2}`)
var euroNumberRegexp = regexp.MustCompile(`[A-Z]{3}[0-9]{3}`)

// ToCarNumber returns searchable car number string
func ToCarNumber(cn string) (string, bool) {
	cn = strings.ToUpper(cn)
	cn = replacer.Replace(cn)
	if matchedCn := standardNumberRegexp.FindString(cn); matchedCn != "" {
		return matchedCn, true
	}
	if matchedCn := euroNumberRegexp.FindString(cn); matchedCn != "" {
		return matchedCn, true
	}
	return cn, false
}

// ToBotResponse retruns text that describes search results
func (co CarOwner) ToBotResponse() string {
	var sb strings.Builder
	sb.WriteString("Я знайшов ☺️")
	sb.WriteString(fmt.Sprintf("\nНомер авто: %s", co.ID))
	if co.Owner != "" {
		sb.WriteString(fmt.Sprintf("\nІм'я контакта: %s", co.Owner))
	}
	sb.WriteString(fmt.Sprintf("\nНомер телефону: %s", co.Phones[0]))
	if len(co.Phones) > 1 {
		sb.WriteString(fmt.Sprintf("\nДодатковый: %s", co.Phones[1]))
	}
	return sb.String()
}
