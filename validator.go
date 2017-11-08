package formvalidator

import (
	"errors"
	"net/url"
)

// global variables
var iso3166countryCodes []string
var iso4217currencyCodes []string
var commonPasswordList []string
var disposableDomains []string
var disposableWildcards []string

func init() {
	iso3166countryCodes = mustFormatCSVBytes(MustAsset("data/iso3166-country-codes.csv"))
	commonPasswordList = mustFormatCSVBytes(MustAsset("data/common-password-list-minlength-8-characters.csv"))
	disposableDomains = mustFormatCSVBytes(MustAsset("data/disposable-email-domains.csv"))
	disposableWildcards = mustFormatCSVBytes(MustAsset("data/disposable-email-domains-wildcards.csv"))
	iso4217currencyCodes = mustFormatCSVBytes(MustAsset("data/iso4217-currency-codes.csv"))
}

/*
	Most rules only validate the first string from the slice in url.Values.
	A few require the entire slice from url.Values.
*/
type Rule interface {
	Validate([]string, map[string]string) (error, []interface{})
}

// Setup all the form validation rules
// Note: case sensitive, if form field is "email" and map field is "Email" the validator will fail!
func RuleChain(rules ...Rule) []Rule {
	var ruleSlice []Rule
	for _, r := range rules {
		ruleSlice = append(ruleSlice, r)
	}

	return ruleSlice
}

type FormValidator struct {
	rules                map[string][]Rule
	errorMessages        map[string]string
	blankFormDataOnError bool
}

// errors
var ErrNilArguments = errors.New("Arguments must be non-nil!")

func New(rules map[string][]Rule) (error, *FormValidator) {

	if rules == nil {
		return ErrNilArguments, nil
	}

	// default error messages for invalid form entries
	// extra error messages (those that are not used in field validation functions) are here for convenience, grouped for translation.
	var errors = map[string]string{
		"account":          "The e-mail or password you entered is incorrect.",
		"alpha_num":        "This field may only contain letters and numbers.",
		"boolean":          "This field must be true or false.",
		"captcha":          "The characters you entered did not match the word verification. Please retry.",
		"city":             "We could not find that city. Please check your spelling.",
		"country_code":     "Please select a valid country.",
		"credit_card":      "Please enter a valid credit card number.",
		"csrf":             "There was an error submitting the form. Please retry.",
		"currency_code":    "Please enter a valid currency code.",
		"date":             "This field must be in a date format (DD-MM-YYYY) [Ex: 31-12-1990]",
		"date_time":        "This field must be in a date-time format (DD-MM-YYYY HH:MM:SS) [Ex: 31-12-1990 14:23:56]",
		"delimiter_min":    "Entries must be at least %d characters long.",
		"delimiter_max":    "Entries cannot be more than %d characters long.",
		"duplicate":        "This field cannot contain duplicate entries.",
		"email":            "Please enter a valid e-mail address.",
		"email_taken":      "That e-mail address is already in use.",
		"float":            "This field must be a floating point number. (Example: -10.50)",
		"float_range":      "This field must be between %f - %f.",
		"in_list":          "Please make a selection.",
		"int_range":        "This field must be between %d - %d.",
		"isbn":             "Please enter a valid ISBN.",
		"json":             "This field must contain valid JSON (Javascript object notation).",
		"latitude":         "Latitude must be between -90.0 degrees and 90.0 degrees.",
		"longitude":        "Longitude must be between -180.0 degrees and 180.0 degrees.",
		"multiple_entries": "This field may only contain one entry.",
		"not_in_list":      "This field contains an invalid entry.",
		"numeric":          "This field must contain enter only numbers.",
		"required":         "This field is required.",
		"slug":             "This field must contain at least one letter or number.",
		"string_matches":   "Fields did not match.",
		"string_max":       "This field cannot be more than %d characters long.",
		"string_min":       "This field must be at least %d characters long.",
		"time":             "This field must be in a time format (HH:MM:SS) [Ex: 14:23:56]",
		"unselected_field": "Please select this field.",
		"utf8_letter_num":  "This field may only contain letters and numbers (Character set: UTF8).",
		"uuid":             "Please enter a valid UUID.",
		"weak_password":    "Please use a stronger password.",
		"web_request_uri":  "Please enter a valid Web URI.",
	}

	return nil, &FormValidator{rules, errors, true}
}

// custom error messages option
func (f *FormValidator) SetErrors(msgs map[string]string) error {
	if msgs == nil {
		return ErrNilArguments
	}

	f.errorMessages = msgs
	return nil
}

// option to not blank the form field if there is an error after validating
func (f *FormValidator) SetBlankOnError(b bool) {
	f.blankFormDataOnError = b
}

func (f *FormValidator) GetErrorMessage(key string) string {
	val, ok := f.errorMessages[key]
	if !ok {
		return ""
	}

	return val
}

/*
	Important: This package is case-sensitive! That means a form field named "email" is different than "Email"

	Loop through the rules and validate each entry
	returns (bool, if the form is valid, map for error messages)
*/
func (f *FormValidator) Validate(form url.Values) (bool, map[string][]FormError) {
	allErrors := make(map[string][]FormError)

	for fieldName, ruleSlice := range f.rules { // loop through map fields, Note: case sensitive, if form field is "email" and map field is "Email" the validator will fail!

		val, _ := form[fieldName] // this will be an empty slice if 'fieldName' does not exist in the map

		var errors []FormError
		for _, r := range ruleSlice { // loop through rule slice
			if err, data := r.Validate(val, f.errorMessages); err != nil {
				errors = appendError(errors, &FormError{err.Error(), data}) // format errors for translation (string separate from extra data)
			}
		}
		allErrors[fieldName] = errors // set the errors for the form entry

		if len(errors) > 0 && f.blankFormDataOnError { // blank the field in original form map if there is an error
			form.Set(fieldName, "")
		}

	}

	if len(allErrors) == 0 {
		return true, allErrors
	}

	return false, allErrors
}
