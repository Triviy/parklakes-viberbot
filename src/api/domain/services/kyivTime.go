package services

import (
	"time"

	"github.com/pkg/errors"
)

const kyiv = "Europe/Kiev"
const kyivFormat = "02.01.2006 15:04:05"

// ToKyivTime converts string value to Kyiv time
func ToKyivTime(s string) (t time.Time, err error) {
	loc, err := time.LoadLocation(kyiv)
	if err != nil {
		return t, errors.Wrap(err, "location loading failed")
	}
	if t, err = time.ParseInLocation(kyivFormat, s, loc); err != nil {
		return t, errors.Wrap(err, "parsing time failed")
	}
	return t, nil
}

// GetKyivTime returns current time.Time in Kyiv
func GetKyivTime() (t time.Time, err error) {
	loc, err := time.LoadLocation(kyiv)
	if err != nil {
		return t, errors.Wrap(err, "loading location failed")
	}
	return time.Now().In(loc), nil
}

// ToKyivFormat formats time to Kyiv time format
func ToKyivFormat(t time.Time) string {
	return t.Format(kyivFormat)
}
