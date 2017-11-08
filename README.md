formvalidator
========

[Go](http://golang.org) package for validating just web forms (url.Values), not structs, maps, or unknown types.

You can write custom error messages and create a custom validation rule by implementing the "Rule" interface.

```go
type Rule interface {
	Validate([]string, map[string]string) (error, []interface{})
}
```

It also includes a couple of functions for rendering radio buttons, checkboxes, and single or multiple selects/dropdowns.

## Installation

```bash
go get -u github.com/dholtzmann/formvalidator
```

## Basic Example

```go
package main

import (
	fv "github.com/dholtzmann/formvalidator"
	"html/template"
	"net/http"
)

func main() {
	// setup HTTP server...
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("./templates/*.html"))

	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	/*
		For file uploads and a form.

		err := r.ParseMultipartForm(1 << 20)
		if err != nil {
			panic(err.Error())
		}
	*/

	rules := map[string][]fv.Rule{
		"Email": fv.RuleChain(fv.Required(), fv.Email(true), fv.StrLen(6, 50)),
		"Age":   fv.RuleChain(fv.Numeric(), fv.IntRange(18, 100)),
		// More validation rules...
	}

	err, validator := fv.New(rules)
	if err != nil {
		panic(err.Error())
	}

	isValid, errors := validator.Validate(r.Form) // or r.PostForm, r.MultipartForm.Value
	if isValid {
		// save to database, parse, ...
	} else {
		// display errors
		tplVars := make(map[string]interface{})
		tplVars["FormErrors"] = errors
		err = tpl.ExecuteTemplate(w, "index.gohtml", tplVars)
	}
}
```

## Validation rules

#### "rules-string.go"
- Required()
- RequiredMultiple()
- Email(allowThrowawayDomains bool)
- AlphaNumeric()
- MinStrLen(min uint32)
- MaxStrLen(max uint32)
- StrLen(min, max uint32)
- UTF8LetterNum()
- StrMatch(comparison string)
- IsJSON()
- WebRequestURI()
- Boolean()

#### "rules-numeric.go"
- Numeric()
- IntRange(min, max int)
- IsFloat64()
- FloatRange(min, max float64)
- CreditCard(allowTestingNumbers bool)
- Latitude()
- Longitude()
- ISBN() [format: Strlen 10 or 13]
- IsUUID(version uint32) [format: UUIDv3, UUIDv4, UUIDv5]
- IsDate() [format: DD-MM-YYYY]
- IsTime() [format: HH:MM:SS]
- IsDateTime() [format: DD-MM-YYYY HH:MM:SS]

#### "rules-list.go"
- CSVEntryStrLen(delimiter rune, minLenValue, maxLenValue uint32)
- CountryCode() [format: ISO-3166, 2 letters]
- InListMultiple(list []string)
- InListSingle(list []string)
- NotInListSingle(list []string)
- CurrencyCode() [format: ISO-4217, 3 letters]
- NotCommonPassword()

See the file 'validator_test.go' for examples of how to use the validation rules.

## Todo
* Write Go documentation in comments. (typical godoc.org format)
* Expand README.md
