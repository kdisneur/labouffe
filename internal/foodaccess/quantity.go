package foodaccess

import (
	"fmt"
	"strings"
)

// QuantityUnit represents the quantity unit
type QuantityUnit int

// Quantity represents an ingredient quantity and can be
// represented in gramm (g) or centiliter (cl)
type Quantity struct {
	Value float64
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

	// QuantityUnitTeaSpoon represents a tea spoon
	QuantityUnitTeaSpoon // cc
	// QuantityUnitTableSpoon represents a table spoon
	QuantityUnitTableSpoon // cs

	// QuantityUnitNoUnit represents a single unit as in one onion
	QuantityUnitNoUnit //
)

func (q Quantity) String() string {
	return fmt.Sprintf("%g%s", q.Value, q.Unit)
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
		var quantity float64

		if _, err := fmt.Sscanf(rawValue, fmt.Sprintf("%%f%s ", allNamedUnits[i]), &quantity); err != nil {
			continue
		}

		q.Value = quantity
		q.Unit = QuantityUnit(i)

		return nil
	}

	return fmt.Errorf("unsupported quantity '%s'. valid format are: %s", rawValue, strings.Join(allNamedUnits, ", "))
}
