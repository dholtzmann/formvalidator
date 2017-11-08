package formvalidator

import (
	"testing"
)

func Test_requiredSingle(t *testing.T) {
	var list = []struct {
		field       []string
		expectation bool
	}{
		{[]string{"true"}, true},
		{[]string{"10"}, true},
		{[]string{"-5.20301"}, true},
		{[]string{"a"}, true},
		{[]string{""}, false},
		{[]string{"a"}, true},
		{[]string{"1"}, true},
		{[]string{"."}, true},
		{[]string{"0"}, true},
		{[]string{"!%^&*"}, true},
		{[]string{"あ"}, true},
		{[]string{"あいうえお"}, true},
		{[]string{"Really long string goes here"}, true},
		{[]string{"a", "b"}, false},
		{[]string{"", ""}, false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = Required()

		if e, _ := rule.Validate(l.field, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("requiredSingle(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_requiredMultiple(t *testing.T) {
	var list = []struct {
		field       []string
		expectation bool
	}{
		{[]string{"true"}, true},
		{[]string{"10"}, true},
		{[]string{"-5.20301"}, true},
		{[]string{"a"}, true},
		{[]string{""}, true},
		{[]string{"1"}, true},
		{[]string{"."}, true},
		{[]string{"0"}, true},
		{[]string{"!%^&*"}, true},
		{[]string{"Really long string goes here"}, true},
		{[]string{"-1", "4", "3810"}, true},
		{[]string{}, false},
		{[]string{"", "", ""}, true},
		{[]string{"a"}, true},
		{[]string{"あ"}, true},
		{[]string{"あ", "お", "え", "う", "い"}, true},
		{[]string{"a", "b", "c"}, true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = RequiredMultiple()

		if e, _ := rule.Validate(l.field, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("requiredMultiple(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_minLength(t *testing.T) {
	var list = []struct {
		field       string
		min         uint32
		expectation bool
	}{
		{"true", 5, false},
		{"10", 3, false},
		{"-5.20301", 10, false},
		{"", 10, true},
		{"a", 2, false},
		{"a", 1, true},
		{"abc", 3, true},
		{"How do you do?", 10, true},
		{"abcdef", 4, true},
		{"あいうえお", 15, true},
		{"あいうえおか", 18, true},
		{"あいうえ", 12, true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = MinStrLen(l.min)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("minLength(%s, %d): Valid[%t]. Expected: %t", l.field, l.min, valid, l.expectation)
		}
	}
}

func Test_maxLength(t *testing.T) {
	var list = []struct {
		field       string
		max         uint32
		expectation bool
	}{
		{"true", 3, false},
		{"10", 1, false},
		{"-5.20301", 7, false},
		{"", 10, true},
		{"a", 0, false},
		{"a", 1, true},
		{"abc", 3, true},
		{"How do you do?", 14, true},
		{"%$#^&*@", 7, true},
		{"あいうえお", 15, true},
		{"あいうえおか", 18, true},
		{"あいうえ", 12, true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = MaxStrLen(l.max)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("maxLength(%s, %d): Valid[%t]. Expected: %t", l.field, l.max, valid, l.expectation)
		}
	}
}

func Test_stringsMatch(t *testing.T) {
	var list = []struct {
		f1          string
		f2          string
		expectation bool
	}{
		{"true", "truE", false},
		{"10", "1.0", false},
		{"-5.20301", "-5.2030l", false},
		{"a", "ä", false},
		//		{"", "b", true}, // use required, validates because one string is length of 0
		//		{"b", "", true}, // use required
		{"", "", true}, // use required
		{"1.2", "1.2", true},
		{"abcdefg", "abcdefg", true},
		{"~!@#$%^&*()_+", "~!@#$%^&*()_+", true},
		{"Horse staple battery!", "Horse staple battery!", true},
		{"abcdefg", "@bcdefg", false},
		{"abcdefg", "abCdefg", false},
		{"thisthat", "this that", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = StrMatch(l.f2)

		if e, _ := rule.Validate([]string{l.f1}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("stringsMatch('%s', '%s'): Valid[%t]. Expected: %t", l.f1, l.f2, valid, l.expectation)
		}
	}
}

func Test_alphaNumeric(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", true},
		{"10", true},
		{"-5.20301", false},
		{"5.20301", false},
		{"520301", true},
		{"a", true},
		{"", true},
		{"abc", true},
		{"números", false},
		{"Não", false},
		{"forneça", false},
		{"é", false},
		{"n'est pas", false},
		{"gültige", false},
		{"zulässig", false},
		{"123lkasdf", true},
		{"1.02", false},
		{"-1234", false},
		{"a=b", false},
		{"How are you?", false},
		{"I am great!", false},
		{"Test.", false},
		{"中文网", false},
		{"a@a.com", false},
		{"ẞ", false},
		{"¾", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
	}

	for _, l := range list {
		valid := false
		var rule Rule = AlphaNumeric()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("alphaNumeric(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_utf8LetterNumeric(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", true},
		{"10", true},
		{"-5.20301", false},
		{"5.20301", false},
		{"520301", true},
		{"a", true},
		{"", true},
		{"abc", true},
		{"números", true},
		{"𐅪", true},
		{"𐅪4", true},
		{"Não", true},
		{"forneça", true},
		{"é", true},
		{"n'est pas", false},
		{"gültige", true},
		{"zulässig", true},
		{"ẞ", true},
		{"123lkasdf", true},
		{"1.02", false},
		{"-1234", false},
		{"a=b", false},
		{"How are you?", false},
		{"I am great!", false},
		{"Test.", false},
		{"中文网", true},
		{"소주", true},
		{"a@a.com", false},
		{"¾", true},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"￩", false},
		{"\x19test\x7F", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = UTF8LetterNum()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("utf8LetterNumeric(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isBoolean(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"", true},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"1", true},
		{"t", true},
		{"T", true},
		{"TRUE", true},
		{"true", true},
		{"True", true},
		{"0", true},
		{"f", true},
		{"F", true},
		{"FALSE", true},
		{"false", true},
		{"False", true},
		{"TRue", false},
		{"FAlse", false},
		{"-1", false},
		{"abcd", false},
		{"0.1", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = Boolean()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isBoolean(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_creditCard(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"0000000000000", false},
		{"a0000000000000", false},
		{"4242424242424242", false},
		{"", true},
		{"4.716461583322103", false},
		{" 4 7 164 61 583 322 10 3", false},
		{"4716 4615 8332 2103", false},
		{"blahblah", false},
		{"-5398228707871527", false},
		{"-5398228707871527-", false},
		{"5398-2287-0787-1527", false},
		{"-38.4716461583322103", false},
		{"5398228707871528", false},
		{"375556917985515", true},
		{"36050234196908", true},
		{"4716461583322103", true},
		{"4716221051885662", true},
		{"4929722653797141", true},
		{"5398228707871527", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = CreditCard(false)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("creditCard(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isWebRequestURI(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"http://a.com", true},
		{"http://something.com/#link", true},
		{"http://something.com/?a=b&c=d&thing=123455&what=that&q=%2F", true},
		{"http://abc.中文网/", true},
		{"http://example.com:8080", true},
		{"http://www.a----thing-that-goes.com:8080", true},
		{"http://127.0.0.1/", true},
		{"http://localhost/", true},
		{"https://gmail.com/inbox", true},
		{"http://user:pass@www.foobar.com/", true},
		{"abc://thing.com", false},
		{"something.", false},
		{".com", false},
		{"ftp://website.com/directory", false},
		{"ftps://website.com/directory", false},
		{"irc://irc.server.org/channel", false},
		{"abc.com", false},
		{"!abc.com", false},
		{"@abc.com", false},
		{"http://somEthING@tHaT.CoM", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = WebRequestURI()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isWebRequestURL(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isJSON(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"", true},
		{"10", true},
		{"-5.20301", true},
		{"a", false},
		{"", true},
		{"123:f00", false},
		{"{\"Name\":\"Alice\",\"Body\":\"Hello\",\"Time\":1294706395881547000}", true},
		{"{}", true},
		{"{\"Key\":{\"Key\":{\"Key\":123}}}", true},
		{"[]", true},
		{"null", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = IsJSON()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isJSON('%s'): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_email(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"[]string{}", false},
		{`[]string{"a"}`, false},
		{`[]string{"2"}`, false},
		{`[]string{"abc", "def"}`, false},
		{"", true},
		{"abc", false},
		{"abc.com", false},
		{"!abc.com", false},
		{"@abc.com", false},
		{"a@a.a", true},
		{"a@Não.b", true},
		{"somEthING@tHaT.CoM", true},
		{"a@中.中文网", true},
		{"文@中.中文网", true},
		{"forneça-here@that.com", true},
		{"gültige@Heiẞe.de", true},
		{"zulässig___abc@that.co.uk", true},
		{"______ẞ-123-test@blahdiblah.com", true},
		{"a.b.c.de.f@a.b.c", true},
		{"a.b.c-----de------.f@a.b.c", true},
		{"a.b.c______def@a.b.c", true},
		{"++++++a+b+++++@a.b.c", true},
		{"a+b+++++@a.b.c", true},
		{"b:a@a.b.c", false},
		{"a{}b@a.b.c", true},
		{"a[]b@a.b.c", false},
		{"a()b@a.b.c", false},
		{"a\"@a.b.c", false},
		{"a'@a.b.c", true},
		{"a\\@a.b.c", false},
		{"a/@a.b.c", true},
		{"a|@a.b.c", true},
		{"a~@a.b.c", true},
		{"a&@a.b.c", true},
		{"a^@a.b.c", true},
		{"a?@a.b.c", true},
		{"123@123.abc", true},
		{"123@123.abc.def", true},
		{"123@123.abc______b", true},
		{"123@123.abc--------a", true},
		{"123@1------23.abc", true},
		{"123@12________3.abc", true},
		{"a!@a.b.c", true},
		{"a#@a.b.c", true},
		{"a$@a.b.c", true},
		{"a%@a.b.c", true},
		{"=a@a.b.c", true},
		{"a=@a.b.c", true},
		{"a;b@a.b.c", false},
		{"a,b@a.com", false},
		{"a@a.b.c", true},
		{"\u0070@that.com", true},  //UTF-8(ASCII): p
		{"this@\u0026.net", false}, //UTF-8(ASCII): &
		{"𐅪@thing.com", false},
		{"𐅪4@gmail.com", false},
		{"¾@place.com", true},
		{"abc123|456@something.com", true},
		{"￩left@right.com", true},
		{"\x19test\x7F@gmail.com", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = Email(true)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("email(%s, true): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}
