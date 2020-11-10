package recipe

import (
	"fmt"
	"strings"
)

// QuantityUnit represents the quantity unit
type QuantityUnit int

// Quantity represents an ingredient quantity and can be
// represented in gramm (g) or centiliter (cl)
type Quantity struct {
	Value int
	Unit  QuantityUnit
}

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=QuantityUnit -linecomment
const (
	// QuantityUnitMilliliter represents milliliter unit
	QuantityUnitMilliliter QuantityUnit = iota // ml
	// QuantityUnitCentiliter represents centiliter unit
	QuantityUnitCentiliter // cl
	// QuantityUnitLiter represents liter unit
	QuantityUnitLiter // l

	// QuantityUnitMilligram represents milligram unit
	QuantityUnitMilligram // mg
	// QuantityUnitGram represents gram unit
	QuantityUnitGram // g
	// QuantityUnitKilogram represents kilogam unit
	QuantityUnitKilogram // kg
)

func (q Quantity) String() string {
	return fmt.Sprintf("%d%s", q.Value, q.Unit)
}

// UnmarshalYAML is the function in charge of unmatshalling the string value to a Go constant
func (q *Quantity) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawValue string
	if err := unmarshal(&rawValue); err != nil {
		return fmt.Errorf("can't parse quantity to string: %v", err)
	}

	var allNamedUnits []string
	for i := 0; i < len(_QuantityUnit_index)-1; i++ {
		allNamedUnits = append(allNamedUnits, _QuantityUnit_name[_QuantityUnit_index[i]:_QuantityUnit_index[i+1]])
	}

	for i := range allNamedUnits {
		var quantity int

		if _, err := fmt.Sscanf(rawValue, fmt.Sprintf("%%d%s", allNamedUnits[i]), &quantity); err != nil {
			continue
		}

		if rawValue != fmt.Sprintf("%d%s", quantity, allNamedUnits[i]) {
			// it makes sure we don't have additional characters.
			// For example: `12clx` would match `12cl` if only relying on Sscanf
			continue
		}

		q.Value = quantity
		q.Unit = QuantityUnit(i)

		return nil
	}

	return fmt.Errorf("unsupported quantity '%s'. valid format are: %s", rawValue, strings.Join(allNamedUnits, ", "))
}
