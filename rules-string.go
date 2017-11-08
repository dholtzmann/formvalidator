package formvalidator

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type requiredSingle struct{}

func Required() Rule {
	return &requiredSingle{}
}

/*
	This requires a slice because it allows the function to test if a user submits multiple values for one field entry which should only have one
*/
func (r *requiredSingle) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {
	// only one entry allowed and string must be at least 1 character long

	if len(fields) == 1 && len(fields[0]) > 0 {
		return nil, nil
	}

	if len(fields) > 1 {
		return errors.New(errorMessages["multiple_entries"]), nil
	}

	return errors.New(errorMessages["required"]), nil
}

// -----------------------

type requiredMultiple struct{}

func RequiredMultiple() Rule {
	return &requiredMultiple{}
}

/*
 */
func (r *requiredMultiple) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {
	// slice must have at least one entry
	if len(fields) > 0 {
		return nil, nil
	}

	return errors.New(errorMessages["required"]), nil
}

// -----------------------

type minStrLen struct {
	min uint32
}

func MinStrLen(min uint32) Rule {
	return &minStrLen{min}
}

/*
 */
func (m *minStrLen) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	if len(field) >= int(m.min) {
		return nil, nil
	}

	return errors.New(errorMessages["string_min"]), []interface{}{m.min}
}

// -----------------------

type maxStrLen struct {
	max uint32
}

func MaxStrLen(max uint32) Rule {
	return &maxStrLen{max}
}

/*
 */
func (m *maxStrLen) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	if len(field) <= int(m.max) {
		return nil, nil
	}

	return errors.New(errorMessages["string_max"]), []interface{}{m.max}
}

// -----------------------

type strLen struct {
	min uint32
	max uint32
}

func StrLen(min, max uint32) Rule {
	return &strLen{min, max}
}

/*
 */
func (s *strLen) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	var mi minStrLen = minStrLen{s.min}
	var ma maxStrLen = maxStrLen{s.max}

	// if there are two errors, return the first encountered error
	if err, data := mi.Validate(fields, errorMessages); err != nil {
		return err, data
	}

	if err, data := ma.Validate(fields, errorMessages); err != nil {
		return err, data
	}

	return nil, nil
}

// -----------------------

type strMatch struct {
	comparison string
}

func StrMatch(comparison string) Rule {
	return &strMatch{comparison}
}

/*
	do strings match exactly? (passwords, ...)
*/
func (m *strMatch) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {
	// if one is blank, it is fine, use required
	/*
		if len(s1) == 0 || len(s2) == 0 {
			return nil, nil
		}
	*/

	field := getFirstKey(fields)

	if field == m.comparison {
		return nil, nil
	}

	return errors.New(errorMessages["string_matches"]), nil
}

// -----------------------

type alphaNumeric struct {
}

func AlphaNumeric() Rule {
	return &alphaNumeric{}
}

/*
 */
func (a *alphaNumeric) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	var alphanum = "^[a-zA-Z0-9]+$"

	if match, _ := regexp.MatchString(alphanum, field); match == true {
		return nil, nil
	}

	return errors.New(errorMessages["alpha_num"]), nil
}

// -----------------------

type utf8LetterNum struct {
}

func UTF8LetterNum() Rule {
	return &utf8LetterNum{}
}

/*
 */
func (u *utf8LetterNum) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	valid := true
	// check all the characters before deciding
	// numbers cannot have decimals or negative sign
	for _, c := range field {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			valid = false
			break
		}
	}

	if valid {
		return nil, nil
	}

	return errors.New(errorMessages["utf8_letter_num"]), nil
}

// -----------------------

type boolean struct {
}

func Boolean() Rule {
	return &boolean{}
}

/*
	From strconv package: It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. Any other value returns an error.
*/
func (b *boolean) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// see function documentation
	if _, err := strconv.ParseBool(field); err == nil {
		return nil, nil
	}

	return errors.New(errorMessages["boolean"]), nil
}

// -----------------------

type creditCard struct {
	allowTestingNumbers bool
}

func CreditCard(allowTestingNumbers bool) Rule {
	return &creditCard{allowTestingNumbers}
}

/*
	Luhn algorithm
	Note: Does not allow testing card numbers! card numbers can have all digits, or be separated by spaces or hypthens
*/
func (c *creditCard) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// change this to remove after 4 digits?
	// remove spaces and hypthens
	//	field = strings.Replace(field, " ", "", -1)
	//	field = strings.Replace(field, "-", "", -1)

	// check for testing card numbers
	// from https://stripe.com/docs/testing
	if c.allowTestingNumbers == false {
		switch field {
		case "4242424242424242",
			"4012888888881881",
			"4000056655665556",
			"5555555555554444",
			"5200828282828210",
			"5105105105105100",
			"378282246310005",
			"371449635398431",
			"6011111111111117",
			"6011000990139424",
			"30569309025904",
			"38520000023237",
			"3530111333300000",
			"3566002020360505",

			// all zeroes pass the luhn algorithm
			"0000000000000",
			"00000000000000",
			"000000000000000",
			"0000000000000000",
			"00000000000000000",
			"000000000000000000",
			"0000000000000000000":
			return errors.New(errorMessages["credit_card"]), nil
		}
	}

	// taken from 'go-credit-card', Luhn algorithm implementation
	var sum int
	var alternate bool

	numberLen := len(field)

	if numberLen < 13 || numberLen > 19 {
		return errors.New(errorMessages["credit_card"]), nil
	}

	for i := numberLen - 1; i > -1; i-- {
		mod, err := strconv.Atoi(string(field[i]))
		if err != nil {
			return errors.New(errorMessages["credit_card"]), nil
		}
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate

		sum += mod
	}

	if sum%10 == 0 {
		return nil, nil
	}

	return errors.New(errorMessages["credit_card"]), nil
}

// -----------------------

type webRequestURI struct {
}

func WebRequestURI() Rule {
	return &webRequestURI{}
}

/*
	Does the URL confirm to RFC 3986 - net/url.ParseRequestURI("...")
	This function expects a website address/URL, not just any valid URL string, it only accepts urls with the protocol "http/https"
*/
func (w *webRequestURI) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// parse raw uri
	u, err := url.ParseRequestURI(field)
	if err != nil {
		return errors.New(errorMessages["web_request_uri"]), nil
	}

	// only allow http/https, not ftp/ftps/mailto/irc/rtmp/etc.
	scheme := strings.ToLower(u.Scheme)
	if !(scheme == "http" || scheme == "https") {
		return errors.New(errorMessages["web_request_uri"]), nil
	}

	return nil, nil
}

// -----------------------

type isJSON struct {
}

func IsJSON() Rule {
	return &isJSON{}
}

/*
	Use encoding/json Unmarshal to test a field
*/
func (i *isJSON) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	var holder json.RawMessage
	if err := json.Unmarshal([]byte(field), &holder); err != nil { // parse json
		return errors.New(errorMessages["json"]), nil
	}

	return nil, nil
}

// -----------------------

type email struct {
	allowDisposableDomains bool
}

func Email(allowDisposableDomains bool) Rule {
	return &email{allowDisposableDomains}
}

/*
	Tests e-mail format and optionally if the e-mail address is a "throw-away" account
*/
func (e *email) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	field := getFirstKey(fields)

	// if blank, it is fine
	if len(field) == 0 {
		return nil, nil
	}

	// regexp taken from 'govalidator'
	var email = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

	if match, _ := regexp.MatchString(email, field); match == true {

		if e.allowDisposableDomains == false {

			// is the domain a throwaway one?

			// extract the domain from the e-mail address
			d := strings.LastIndex(field, "@")
			domain := field[d+1:]

			// it is in the list, bad e-mail address
			if inSlice(disposableDomains, domain) {
				return errors.New(errorMessages["email"]), nil
			}

			// wildcard list
			if inSlice(disposableWildcards, domain) {
				return errors.New(errorMessages["email"]), nil
			}
		}

		return nil, nil
	}

	return errors.New(errorMessages["email"]), nil
}
