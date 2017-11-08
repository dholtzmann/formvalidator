package formvalidator

import (
	"net/url"
	"testing"
)

func TestValidate(t *testing.T) {
	form := url.Values{}
	form.Set("FirstName", "Henry")
	form.Set("LastName", "Smith")
	form.Set("Username", "hensmith55")
	form.Set("Email", "henrysmith@website.com")
	form.Set("Password", "I am Great!")
	form.Set("ConfirmPassword", "I am Great!")
	form.Set("Age", "55")
	form.Set("Amount", "22.50")
	form.Set("CreditCard", "4716461583322103")
	form.Set("Latitude", "45.3941")
	form.Set("Longitude", "134.03843")
	form.Set("Website", "https://henrysmithisgreat.com/")
	form.Set("Date", "17-09-1974")
	form.Set("Time", "03:37:19")
	form.Set("DateTime", "29-04-1784 07:19:55")
	form.Set("Slug", "Test123%")
	form.Set("DelimiterList", ",,,,one,two,three,,,,,four,five,,,,")
	form.Set("AgreeToTerms", "true")
	form.Set("Country", "gb")
	form.Add("FavoriteColors", "orange")
	form.Add("FavoriteColors", "black")
	form.Add("FavoriteAnimal", "dogs")
	form.Set("LeastFavoriteAnimal", "mosquitoes")
	form.Set("Currency", "usd")
	form.Set("ISBN", "3 401 01319 X")
	form.Set("UUID4", "57b73598-8764-4ad0-a76a-679bb6640eb1")
	form.Set("JSON", `{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)

	// this should validate without a problem
	rules := map[string][]Rule{
		"FirstName":           RuleChain(Required(), AlphaNumeric(), StrLen(2, 50)),
		"LastName":            RuleChain(Required(), AlphaNumeric(), StrLen(4, 50)),
		"Username":            RuleChain(Required(), UTF8LetterNum(), StrLen(2, 30)),
		"Email":               RuleChain(Required(), Email(true), StrLen(6, 50)),
		"Password":            RuleChain(Required(), StrMatch(form.Get("ConfirmPassword")), StrLen(8, 500), NotCommonPassword()),
		"ConfirmPassword":     RuleChain(Required()),
		"Age":                 RuleChain(Numeric(), IntRange(18, 100)),
		"Amount":              RuleChain(Required(), IsFloat64(), FloatRange(0.01, 100.0)),
		"CreditCard":          RuleChain(Required(), CreditCard(false)),
		"Latitude":            RuleChain(Required(), Latitude()),
		"Longitude":           RuleChain(Required(), Longitude()),
		"Website":             RuleChain(WebRequestURI()),
		"Date":                RuleChain(IsDate()),
		"Time":                RuleChain(IsTime()),
		"DateTime":            RuleChain(IsDateTime()),
		"DelimiterList":       RuleChain(CSVEntryStrLen(',', 2, 5)),
		"AgreeToTerms":        RuleChain(Required(), Boolean()),
		"Country":             RuleChain(Required(), CountryCode()),
		"FavoriteColors":      RuleChain(RequiredMultiple(), InListMultiple([]string{"blue", "green", "red", "yellow", "purple", "white", "black", "brown", "orange"})),
		"FavoriteAnimal":      RuleChain(Required(), InListSingle([]string{"dogs", "cats", "hamsters", "lizards", "birds", "fish"})),
		"LeastFavoriteAnimal": RuleChain(NotInListSingle([]string{"dogs", "cats", "hamsters", "lizards", "birds", "fish"})),
		"Currency":            RuleChain(CurrencyCode()),
		"ISBN":                RuleChain(ISBN()),
		"UUID4":               RuleChain(IsUUID(4)),
		"JSON":                RuleChain(IsJSON()),
	}

	err, validator := New(rules)
	if err != nil {
		t.Errorf("Error making a new form validator type! %s", err.Error())
	}

	_, errors := validator.Validate(form)

	if len(errors) > 0 {
		for field, val := range errors {
			for _, e := range val {
				t.Errorf("TestValidate(): %s: %s", field, e.Error())
			}
		}
	}

	form.Set("FirstName", "")
	form.Set("LastName", "!@#")
	// test the options
	validator.SetBlankOnError(false)

	errorMsgs := map[string]string{
		"required":  "What do you know Joe?",
		"alpha_num": "Crazy man!",
	}
	validator.SetErrors(errorMsgs)
	isValid, errors := validator.Validate(form)

	if isValid != false {
		t.Errorf("TestValidate(): validation should have failed!: %b", isValid)
	}

	if errors["FirstName"][0].Error() != "What do you know Joe?" {
		t.Errorf("TestValidate(): Setting custom error messages failed!")
	}

	if len(errors["LastName"]) != 2 {
		t.Errorf("TestValidate(): LastName should have two errors! [%v]", errors)
	}

	if form.Get("FirstName") != "" {
		t.Errorf("TestValidate(): FirstName should be blank in the form! [%s]", form.Get("FirstName"))
	}

	if form.Get("LastName") != "!@#" { // flag SetBlankOnError(false) should prevent clearing the form field
		t.Errorf("TestValidate(): LastName should not be blank in the form! [%s]", form.Get("LastName"))
	}
}
