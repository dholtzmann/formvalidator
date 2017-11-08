package formvalidator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type numeric struct {
}

func Numeric() Rule {
	return &numeric{}
}

/*
	tests for digits [0-9]
	return error - FormError with message and extra data if revelant
*/
func (n *numeric) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	numeric := "^[0-9]+$"

	if match, _ := regexp.MatchString(numeric, field); match == true {
		return nil, nil
	}

	return errors.New(errorMessages["numeric"]), nil
}

// -----------------------

type isFloat64 struct {
}

func IsFloat64() Rule {
	return &isFloat64{}
}

/*
	tests for float64 (positive or negative and decimal points)
	return error - FormError with message and extra data if revelant
*/
func (f *isFloat64) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// convert to float64
	if _, e := strconv.ParseFloat(field, 64); e == nil {
		return nil, nil
	}

	return errors.New(errorMessages["float"]), nil
}

// -----------------------

type intRange struct {
	min int
	max int
}

func IntRange(min, max int) Rule {
	return &intRange{min, max}
}

/*
	is some number between MIN and MAX? (integers only)
	return error - FormError with message and extra data if revelant
*/
func (r *intRange) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	i, e := strconv.Atoi(field)
	if e == nil && i >= r.min && i <= r.max {
		return nil, nil
	}

	return errors.New(errorMessages["int_range"]), []interface{}{r.min, r.max}
}

// -----------------------

type floatRange struct {
	min float64
	max float64
}

func FloatRange(min, max float64) Rule {
	return &floatRange{min, max}
}

/*
	is some number between MIN and MAX? (float64)
	return error - FormError with message and extra data if revelant
*/
func (r *floatRange) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	i, e := strconv.ParseFloat(field, 64)
	if e == nil && i >= r.min && i <= r.max {
		return nil, nil
	}

	return errors.New(errorMessages["float_range"]), []interface{}{r.min, r.max}
}

// -----------------------

type latitude struct {
}

func Latitude() Rule {
	return &latitude{}
}

/*
	return error - FormError with message and extra data if revelant
*/
func (l *latitude) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// regexp taken from 'govalidator'
	var latitud = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"

	if match, _ := regexp.MatchString(latitud, field); match == true {
		return nil, nil
	}

	return errors.New(errorMessages["latitude"]), nil
}

// -----------------------

type longitude struct {
}

func Longitude() Rule {
	return &longitude{}
}

/*
	return error - FormError with message and extra data if revelant
*/
func (l *longitude) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// regexp taken from 'govalidator'
	var longitud = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"

	if match, _ := regexp.MatchString(longitud, field); match == true {
		return nil, nil
	}

	return errors.New(errorMessages["longitude"]), nil
}

// -----------------------

type isDate struct {
}

func IsDate() Rule {
	return &isDate{}
}

/*
	Is a string a date? (format: DD-MM-YYYY)
	return error - FormError with message and extra data if revelant
*/
func (d *isDate) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// DD-MM-YYYY
	if _, err := time.Parse("02-01-2006", field); err == nil {
		return nil, nil
	}

	return errors.New(errorMessages["date"]), nil
}

// -----------------------

type isTime struct {
}

func IsTime() Rule {
	return &isTime{}
}

/*
	Is a string a time? (format: HH:MM:SS)
	return error - FormError with message and extra data if revelant
*/
func (t *isTime) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// HH:MM:SS
	if _, err := time.Parse("15:04:05", field); err == nil {
		return nil, nil
	}

	return errors.New(errorMessages["time"]), nil
}

// -----------------------

type isDateTime struct {
}

func IsDateTime() Rule {
	return &isDateTime{}
}

/*
	Is a string a datetime? (format: DD-MM-YYYY HH:MM:SS) These are used for timestamps.
	return error - FormError with message and extra data if revelant
*/
func (i *isDateTime) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// DD-MM-YYYY HH:MM:SS
	if _, err := time.Parse("02-01-2006 15:04:05", field); err == nil {
		return nil, nil
	}

	return errors.New(errorMessages["date_time"]), nil
}

// -----------------------

type isUUID struct {
	regex string
}

func IsUUID(version uint32) Rule {
	var regex string
	switch version { // regexp taken from 'govalidator'
	case 3:
		regex = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	case 4:
		regex = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	case 5:
		regex = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"

	default:
		regex = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	}

	return &isUUID{regex}
}

/*
	UUID - universally unique id, used in databases, software, ...
	return error - FormError with message and extra data if revelant
*/
func (i *isUUID) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	if match, _ := regexp.MatchString(i.regex, field); match == true {
		return nil, nil
	}

	return errors.New(errorMessages["uuid"]), nil
}

// -----------------------

type isbn struct {
}

func ISBN() Rule {
	return &isbn{}
}

/*
 */
func (i *isbn) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	r := strings.NewReplacer(" ", "", "-", "") // remove whitespace and hypthens
	field = r.Replace(field)

	if len(field) == 10 {
		rule := ISBN10()
		return rule.Validate(fields, errorMessages)
	} else if len(field) == 13 {
		rule := ISBN13()
		return rule.Validate(fields, errorMessages)
	}

	return errors.New(errorMessages["isbn"]), nil
}

// -----------------------

type isbn10 struct {
}

func ISBN10() Rule {
	return &isbn10{}
}

/*
 */
func (i *isbn10) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	r := strings.NewReplacer(" ", "", "-", "") // remove whitespace and hypthens
	field = r.Replace(field)

	var ISBN10 string = "^(?:[0-9]{9}X|[0-9]{10})$"

	if match, _ := regexp.MatchString(ISBN10, field); match == false {
		return errors.New(errorMessages["isbn"]), nil
	}

	// run checksum algorithm
	var checksum, j uint32
	for j = 0; j < 9; j++ {
		checksum += (j + 1) * uint32(field[j]-'0')
	}
	if field[9] == 'X' {
		checksum += 10 * 10
	} else {
		checksum += 10 * uint32(field[9]-'0')
	}
	if checksum%11 == 0 {
		return nil, nil
	}

	return errors.New(errorMessages["isbn"]), nil
}

// -----------------------

type isbn13 struct {
}

func ISBN13() Rule {
	return &isbn13{}
}

/*
 */
func (i *isbn13) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	r := strings.NewReplacer(" ", "", "-", "") // remove whitespace and hypthens
	field = r.Replace(field)

	var ISBN13 string = "^(?:[0-9]{13})$"

	if match, _ := regexp.MatchString(ISBN13, field); match == false {
		return errors.New(errorMessages["isbn"]), nil
	}

	// run checksum algorithm
	var j, checksum uint32
	factor := []uint32{1, 3}
	for j = 0; j < 12; j++ {
		checksum += factor[j%2] * uint32(field[j]-'0')
	}
	if (uint32(field[12]-'0'))-((10-(checksum%10))%10) == 0 {
		return nil, nil
	}

	return errors.New(errorMessages["isbn"]), nil
}
