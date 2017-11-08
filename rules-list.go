package formvalidator

import (
	"encoding/csv"
	"errors"
	"strings"
)

type inListSingle struct {
	list []string
}

func InListSingle(list []string) Rule {
	if list == nil {
		list = []string{}
	}
	return &inListSingle{list}
}

/*
	This is used in radio buttons and drop-downs (single entry)

	fields []string - only one item is permitted, but with argument being a slice, it allows testing for multiple values
	list []string - set by the programmer so it must be a slice

	return error - FormError with message and extra data if revelant
*/
func (i *inListSingle) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	if len(fields) == 0 || (len(fields) == 1 && len(fields[0]) == 0) { // if blank, it is fine, use required_single
		return nil, nil
	}

	if len(fields) > 1 { // more than one entry, return error
		return errors.New(errorMessages["multiple_entries"]), nil
	}

	// Is the submitted entry in the allowed list?
	if inSlice(i.list, fields[0]) {
		return nil, nil
	}

	return errors.New(errorMessages["in_list"]), nil
}

// -----------------------

type inListMultiple struct {
	list []string
}

func InListMultiple(list []string) Rule {
	if list == nil {
		list = []string{}
	}
	return &inListMultiple{list}
}

/*
	This is used in checkboxes and drop-down multi-select

	field []string - must be a slice for multiple entries
	list []string - set by the programmer so it must be a slice

	return error - FormError with message and extra data if revelant
*/
func (i *inListMultiple) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	// if blank, it is fine, use required_multiple
	if len(fields) == 0 {
		return nil, nil
	}

	// check for duplicate entries
	if hasDuplicates(fields) {
		return errors.New(errorMessages["duplicate"]), nil
	}

	if diff := sliceDiff(fields, i.list); diff == nil {
		return nil, nil
	}

	return errors.New(errorMessages["in_list"]), nil
}

// -----------------------

type notInListSingle struct {
	list []string
}

func NotInListSingle(list []string) Rule {
	if list == nil {
		list = []string{}
	}
	return &notInListSingle{list}
}

/*
	field string - only one item is permitted
	list []string - set by the programmer so it must be a slice

	return error - FormError with message and extra data if revelant
*/
func (i *notInListSingle) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine, use required_single
	if len(field) == 0 {
		return nil, nil
	}

	// Is the submitted entry in the allowed list?
	if !inSlice(i.list, field) {
		return nil, nil
	}

	return errors.New(errorMessages["not_in_list"]), nil
}

// -----------------------

type countryCode struct {
}

func CountryCode() Rule {
	return &countryCode{}
}

/*
	Expects lowercase two letter country codes
	return error - FormError with message and extra data if revelant
*/
func (c *countryCode) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// check the list
	if inSlice(iso3166countryCodes, field) {
		return nil, nil
	}

	return errors.New(errorMessages["country_code"]), nil
}

// -----------------------

type csvEntryStrLen struct {
	delimiter rune
	min       uint32
	max       uint32
}

/*
	csv = comma separated values

	This splits an string based on some delimiter (default is a comma, but can be different), and checks minlen and maxlen of each entry.

	return error - FormError with message and extra data if revelant
*/
func CSVEntryStrLen(delimiter rune, min, max uint32) Rule {
	switch delimiter {
	case ',':
	case '|':

	default:
		delimiter = ','
	}

	return &csvEntryStrLen{delimiter, min, max}
}

func (c *csvEntryStrLen) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {
	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	r := csv.NewReader(strings.NewReader(field))
	r.Comma = c.delimiter

	entries, err := r.ReadAll()
	if err != nil {
		return errors.New(errorMessages["delimiter_min"]), []interface{}{c.min}
	}

	if len(entries) == 0 {
		return errors.New(errorMessages["delimiter_min"]), []interface{}{c.min}
	}

	for _, e := range entries[0] {
		if len(e) == 0 {
			continue
		}

		if len(e) < int(c.min) {
			return errors.New(errorMessages["delimiter_min"]), []interface{}{c.min}
		}

		if len(e) > int(c.max) {
			return errors.New(errorMessages["delimiter_max"]), []interface{}{c.max}
		}
	}

	return nil, nil
}

// -----------------------

type notCommonPassword struct {
}

func NotCommonPassword() Rule {
	return &notCommonPassword{}
}

/*
	return error - FormError with message and extra data if revelant
*/
func (n *notCommonPassword) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// Is the submitted entry in the allowed list?
	if !inSlice(commonPasswordList, field) {
		return nil, nil
	}

	return errors.New(errorMessages["weak_password"]), nil
}

// -----------------------

type currencyCode struct {
}

func CurrencyCode() Rule {
	return &currencyCode{}
}

/*
	Expects lowercase two letter country codes
	return error - FormError with message and extra data if revelant
*/
func (c *currencyCode) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// check the list
	if inSlice(iso4217currencyCodes, field) {
		return nil, nil
	}

	return errors.New(errorMessages["currency_code"]), nil
}
