package models_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triviy/parklakes-viberbot/domain/models"
)

type carResult struct {
	carNumber string
	matched   bool
}

var cars = map[string]carResult{
	"АИ1234АР":            {"AI1234AP", true},
	"АА 12 34 ТХ":         {"AA1234TX", true},
	"АА1234ЕВ або СА":     {"AA1234EB", true},
	"JCK 123":             {"JCK123", true},
	"MEKAN":               {"MEKAN", false},
	"аа1234кі":            {"AA1234KI", true},
	"12345кв":             {"12345KB", true},
	"А123оу123":           {"A123OY123", false},
	"аа1234то":            {"AA1234TO", true},
	"123-12AX":            {"12312AX", true},
	"Volkswagen АА1234СТ": {"AA1234CT", true},
	"AI12-34HE":           {"AI1234HE", true},
	"BLD1234":             {"BLD1234", false},
	"":                    {"", false},
}

func TestToCarNumber(t *testing.T) {
	for v, expected := range cars {
		actualCarNumber, actualMatched := models.ToCarNumber(v)
		passed := assert.Equal(t, expected.carNumber, actualCarNumber) &&
			assert.Equal(t, expected.matched, actualMatched, fmt.Sprintf("Matched should be `%v` for `%s`", actualMatched, actualCarNumber))
		if passed {
			t.Logf("ToCarNumber() passed for %s", expected.carNumber)
		}
	}
}
