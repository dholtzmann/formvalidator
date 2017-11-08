package formvalidator

import (
	"testing"
)

func Test_inListSingle(t *testing.T) {
	var list = []struct {
		entry       []string
		allowed     []string
		expectation bool
	}{
		{[]string{"truE"}, []string{"true"}, false},
		{[]string{"1.0"}, []string{"10"}, false},
		{[]string{"5.20301"}, []string{"-5.20301"}, false},
		{[]string{"a"}, []string{"ä"}, false},
		{[]string{""}, []string{}, true},
		{[]string{"1.2"}, []string{"1.2", "3.4", "5.6"}, true},
		{[]string{"11"}, []string{"00", "11", "22", "33", "44", "55", "66", "77", "88", "99"}, true},
		{[]string{"duDe"}, []string{"dude"}, false},
		{[]string{"dude!"}, []string{"dude"}, false},
		{[]string{"@ttention"}, []string{"attention"}, false},
		{[]string{"~!@#$%^&*()_+"}, []string{"~!@#$%^&*()_+"}, true},
		{[]string{"ä"}, []string{"a", "b"}, false},
		{[]string{"a", ""}, []string{"a", "b"}, false},
		{[]string{"", ""}, []string{"a", "b"}, false},
		{[]string{"c"}, []string{"a", "b"}, false},
		{[]string{"1.355"}, []string{"1.35"}, false},
		{[]string{"O"}, []string{"0"}, false},
	}
	for _, l := range list {
		valid := false
		var rule Rule = InListSingle(l.allowed)

		if e, _ := rule.Validate(l.entry, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("inListSingle('%v', '%v'): Valid[%t]. Expected: %t", l.entry, l.allowed, valid, l.expectation)
		}
	}
}

func Test_inListMultiple(t *testing.T) {
	var list = []struct {
		entry       []string
		allowed     []string
		expectation bool
	}{
		{[]string{"True"}, []string{"true"}, false},
		{[]string{"1O"}, []string{"10"}, false},
		{[]string{"5.20301"}, []string{"-5.20301"}, false},
		{[]string{"a"}, []string{"ä"}, false},
		{[]string{""}, []string{}, false},
		{[]string{}, []string{}, true},
		{[]string{""}, []string{""}, true},
		{[]string{}, []string{""}, true},
		{[]string{"1.2"}, []string{"1.2", "3.4", "5.6"}, true},
		{[]string{"1.2", "1.2"}, []string{"1.2", "3.4", "5.6"}, false},
		{[]string{""}, []string{"00", "11", "22", "33", "44", "55", "66", "77", "88", "99"}, false},
		{[]string{"", ""}, []string{"00", "11", "22", "33", "44", "55", "66", "77", "88", "99"}, false},
		{[]string{"11"}, []string{"00", "11", "22", "33", "44", "55", "66", "77", "88", "99"}, true},
		{[]string{"11", "33", "55", "77", "99"}, []string{"00", "11", "22", "33", "44", "55", "66", "77", "88", "99"}, true},
		{[]string{"duDe"}, []string{"dude", "DUDE"}, false},
		{[]string{"dude!"}, []string{"dude", "DUDE"}, false},
		{[]string{"@ttention"}, []string{"attention"}, false},
		{[]string{"~!@#$%^&*()_+"}, []string{"~!@#$%^&*()_+"}, true},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"c", "a"}, []string{"a", "b"}, false},
		{[]string{"1.35"}, []string{"l.35"}, false},
		{[]string{"0", "1", "2", "3"}, []string{"0", "1", "2", "3", "4"}, true},
	}
	for _, l := range list {
		valid := false
		var rule Rule = InListMultiple(l.allowed)

		if e, _ := rule.Validate(l.entry, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("inListMultiple('%v', '%v'): Valid[%t]. Expected: %t", l.entry, l.allowed, valid, l.expectation)
		}
	}
}

func Test_countryCode(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"u", false},
		{"uS", false},
		{"US", false},
		{"us", true},
		{"usa", false},
		{"United States of America", false},
		{"ab", false},
		{"gb", true},
	}
	for _, l := range list {
		valid := false
		var rule Rule = CountryCode()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		} else {
			//t.Logf("CountryCode(%s): %s", l.field, e)
		}

		if l.expectation != valid {
			t.Errorf("countryCode(%v): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_notInListSingle(t *testing.T) {
	var list = []struct {
		entry       string
		allowed     []string
		expectation bool
	}{
		{"true", []string{"true"}, false},
		{"True", []string{"true"}, true},
		{"1O", []string{"10"}, true},
		{"10", []string{"10"}, false},
		{"-5.20301", []string{"-5.20301"}, false},
		{"ä", []string{"a"}, true},
		{"a", []string{"a"}, false},
		{"", []string{}, true},
		{"something.com", []string{"something.co"}, true},
		{"something.com", []string{"something.coM"}, true},
		{"something.com", []string{"something.com"}, false},
		{"1.20", []string{"1.2", "3.4", "5.6"}, true},
		{"duDe", []string{"dude"}, true},
		{"dude!", []string{"dude"}, true},
		{"@ttention", []string{"attention"}, true},
		{"~!@#$%^&*()_+", []string{"~!@#$%^&*()_+"}, false},
		{"a", []string{"a", "b"}, false},
		{"ä", []string{"a", "b"}, true},
		{"c", []string{"a", "b"}, true},
		{"1.35", []string{"1.35"}, false},
		{"0", []string{"0"}, false},
	}
	for _, l := range list {
		valid := false
		var rule Rule = NotInListSingle(l.allowed)

		if e, _ := rule.Validate([]string{l.entry}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("notInListSingle('%v', '%v'): Valid[%t]. Expected: %t", l.entry, l.allowed, valid, l.expectation)
		}
	}
}

func Test_isCSVEntryMinStrLen(t *testing.T) {
	var list = []struct {
		field       string
		min         uint32
		expectation bool
	}{
		{"true", 5, false},
		{"10", 3, false},
		{"-5.20301", 9, false},
		{"a", 2, false},
		{"", 10, true},
		{"a", 1, true},
		{"abc", 3, true},
		{"How do you do?", 10, true},
		{"abcdef", 4, true},
		{"How do you do?, testing one two three, here", 4, true},
		{"How do you (@*#$&(@#$do?, testing !@#*$&@*()#*$)one two three, here", 4, true},
		{"one,two,", 4, false},
		{",,,one,two,,,,,", 3, true},
		{",,,あ,い,,,,,", 3, true},
		{"あいうえお", 15, true},
		{"あいうえおか", 18, true},
		{"あいうえ", 12, true},
	}

	// test minimum length
	for _, l := range list {
		valid := false
		var rule Rule = CSVEntryStrLen(',', l.min, 50)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("Testing Min. CSVEntryStrLen(%v, ',', %d, 50): Valid[%t]. Expected: %t", l.field, l.min, valid, l.expectation)
		}
	}

	// test other delimiter
	var rule Rule = CSVEntryStrLen('|', 4, 50)

	if e, _ := rule.Validate([]string{"How do you do?| testing one two three| here"}, make(map[string]string)); e != nil {
		t.Errorf("Testing Min. CSVEntryStrLen('How do you do?| testing one two three| here', '|', 4, 50): Valid[false]. Expected: true\n")
	}
}

func Test_isCSVEntryMaxStrLen(t *testing.T) {
	var list = []struct {
		field       string
		max         uint32
		expectation bool
	}{
		{"true", 1, false},
		{"10", 1, false},
		{"-5.20301", 1, false},
		{"a", 0, false},
		{"", 10, true},
		{"a", 1, true},
		{"abc", 3, true},
		{"How do you do?", 14, true},
		{"abcdef", 6, true},
		{"How do you do?, testing one two three, here", 22, true},
		{"How do you (@*#$&(@#$do?, testing !@#*$&@*()#*$)one two three, here", 36, true},
		{"one,two,", 2, false},
		{",,,one,two,,,,,", 3, true},
		{",,,あ,い,,,,,", 3, true},
		{"あいうえお", 15, true},
		{"あいうえおか", 18, true},
		{"あいうえ", 12, true},
	}

	// test maximum length
	for _, l := range list {
		valid := false
		var rule Rule = CSVEntryStrLen(',', 1, l.max)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("Testing Max. CSVEntryStrLen(%v, ',', 1, %d): Valid[%t]. Expected: %t", l.field, l.max, valid, l.expectation)
		}
	}

	// test other delimiter
	var rule Rule = CSVEntryStrLen('|', 1, 22)

	if e, _ := rule.Validate([]string{"How do you do?| testing one two three| here"}, make(map[string]string)); e != nil {
		t.Errorf("Testing Min. CSVEntryStrLen('How do you do?| testing one two three| here', '|', 4, 50): Valid[false]. Expected: true\n")
	}
}

func Test_notInCommonPasswordList(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", true},
		{"10", true},
		{"-5.20301", true},
		{"a", true},
		{"1234567", true}, // less than 8 characters
		{"", true},
		{"playstation3", false},
		{"1234567890", false},
		{"qwertyuiop", false},
		{"1111111111", false},
		{"password123", false},
		{"0123456789", false},
		{"basketball", false},
		{"qazwsxedcrfv", false},
		{"0000000000", false},
		{"5555555555", false},
		{"qwertyqwerty", false},
		{"leavemealone", false},
		{"harrypotter", false},
		{"mypassword", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = NotCommonPassword()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("notInCommonPasswordList(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_currencyCode(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"u", false},
		{"uSd", false},
		{"USd", false},
		{"usd", true},
		{"usa", false},
		{"United States of America", false},
		{"abc", false},
		{"cad", true},
		{"dkk", true},
	}
	for _, l := range list {
		valid := false
		var rule Rule = CurrencyCode()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("currencyCode(%v): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}
