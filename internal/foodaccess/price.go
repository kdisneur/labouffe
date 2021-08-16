package foodaccess

import (
	"fmt"
	"strings"
)

// Price represents the price range of a recipe
type Price int

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=Price -linecomment
const (
	// PriceCheap represents a recipe that is cheap
	PriceCheap Price = iota // économique
	// PriceAffordable represents a recipe that is not cheap but not too expensive
	PriceAffordable // abordable
	// PriceExpensive represents a recipe that is expensive
	PriceExpensive // cher
)

// AllPrices returns the full list of prices
func AllPrices() []Price {
	var prices []Price
	for i := 0; i < len(_Price_index)-1; i++ {
		prices = append(prices, Price(i))
	}

	return prices
}

// UnmarshalYAML is the function in charge of unmarshalling the string value to a Go constant
func (p *Price) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var namedPrice string
	if err := unmarshal(&namedPrice); err != nil {
		return fmt.Errorf("can't unmarshal namedprice, expected a string: %v", err)
	}

	var allNamedPrices []string
	for i := 0; i < len(_Price_index)-1; i++ {
		allNamedPrices = append(allNamedPrices, _Price_name[_Price_index[i]:_Price_index[i+1]])
	}

	for i := range allNamedPrices {
		if allNamedPrices[i] == strings.ToLower(namedPrice) {
			*p = Price(i)
			return nil
		}
	}

	return fmt.Errorf("unsupported price '%s'. valid prices are: %s", namedPrice, strings.Join(allNamedPrices, ", "))
}
