package formvalidator

import (
	"testing"
)

func Test_numeric(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", true},
		{"-5.20301", false},
		{"5.20301", false},
		{"520301", true},
		{"a", false},
		{"", true},
		{"a", false},
		{"1", true},
		{"8792489507", true},
		{".", false},
		{"0", true},
		{"!%^&*", false},
		{"0a", false},
		{"z9", false},
		{"-1", false},
		{"3.14", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = Numeric()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("numeric(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isFloat64(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", true},
		{"-5.20301", true},
		{"5.20301", true},
		{"520301", true},
		{"-520301", true},
		{"a", false},
		{"", true},
		{"a", false},
		{"1", true},
		{"8792489507", true},
		{".", false},
		{"0", true},
		{"!%^&*", false},
		{"0a", false},
		{"z9", false},
		{"-1", true},
		{"3.14", true},
		{"-3821.1432", true},
		{"83291234.1432", true},
		{"12.92193048123", true},
		{"-0.22250738585072011e-307", true},
		{"+0.22250738585072011e-307", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = IsFloat64()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isFloat64(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_intRange(t *testing.T) {
	var list = []struct {
		field       string
		min         int
		max         int
		expectation bool
	}{
		{"true", 1, 2, false},
		{"10", 1, 2, false},
		{"-5.20301", 1, 2, false},
		{"a", 1, 2, false},
		{"", 10, 20, true},
		{"1.2", 1, 2, false},
		{"-0.1", 3, 4, false},
		{"-5", -20, -1, true},
		{"5", 5, 6, true},
		{"0", -1, 1, true},
		{"1234567890", 1000000000, 10000000000, true},
		{"-1234567890", -10000000000, -1000000000, true},
		{"-2147483648", -2147483648, 0, true}, //Signed 32 Bit Min Int
		{"2147483647", 0, 2147483647, true},   //Signed 32 Bit Max Int
		{"-2147483649", -2147483649, 0, true}, //Signed 32 Bit Min Int - 1
		{"2147483648", 0, 2147483648, true},   //Signed 32 Bit Max Int + 1

		{"4294967295", 0, 4294967295, true}, //Unsigned 32 Bit Max Int
		{"4294967296", 0, 4294967296, true}, //Unsigned 32 Bit Max Int + 1

		{"-9223372036854775808", -9223372036854775808, 0, true}, //Signed 64 Bit Min Int
		{"9223372036854775807", 0, 9223372036854775807, true},   //Signed 64 Bit Max Int

	}

	for _, l := range list {
		valid := false
		var rule Rule = IntRange(l.min, l.max)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("intRange(%s, min:%d, max:%d): Valid[%t]. Expected: %t", l.field, l.min, l.max, valid, l.expectation)
		}
	}
}

func Test_floatRange(t *testing.T) {
	var list = []struct {
		field       string
		min         float64
		max         float64
		expectation bool
	}{
		{"true", 1, 2, false},
		{"10", 1, 2, false},
		{"-5.20301", -10.0, -1.25, true},
		{"a", 1, 2, false},
		{"", 10, 20, true},
		{"1.2", 1, 2.5, true},
		{"-0.1234567", -1, 0, true},
		{"5.475390123", 5, 6, true},
		{"0.28934023", -1, 1, true},
		{"12345.12345", 10000, 20000.5, true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = FloatRange(l.min, l.max)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("floatRange(%s, min:%f, max:%f): Valid[%t]. Expected: %t", l.field, l.min, l.max, valid, l.expectation)
		}
	}
}

func Test_latitude(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", true},
		{"-5.20301", true},
		{"a", false},
		{"", true},
		{"-45.112023", true},
		{"71.1203", true},
		{"102.0123", false},
		{"+90.0", true},
		{"90.0", true},
		{"90.1", false},
		{"-90.0", true},
		{"-90.1", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = Latitude()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("latitude(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_longitude(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", true},
		{"-5.20301", true},
		{"a", false},
		{"", true},
		{"-45.112023", true},
		{"71.1203", true},
		{"+71.1203", true},
		{"201.0123", false},
		{"180.0", true},
		{"180.1", false},
		{"-180.0", true},
		{"-180.1", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = Longitude()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("longitude(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

// date and time functions
// format: DD-MM-YYYY
func Test_isDate(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"Thursday", false},
		{"", true},
		{"29-02-2000", true},
		{"29-002-2000", false},
		{"29-02-1999", false},
		{"14-06-2000", true},
		{"14_06_2000", false},
		{"14.06.2000", false},
		{"14/06/2000", false},
		{"14\\06\\2000", false},
		{"32-12-2000", false},
		{"24-13-2000", false},
		{"Monday 22-05-2017", false},
		{"Mon 22-05-2017", false},
		{"Monday 22-May-2017", false},
		{"2000-12-31", false},
		{"2000-31-12", false},
		{"2000/31/12", false},
		{"2000/12/31", false},
		{"99-31-12", false},
		{"24/04/2000", false},
		{"24 04 2000", false},
		{"24 April 2000", false},
		{"24, 04, 2000", false},
		{"24,04,2000", false},
		{"24+04+2000", false},
		{"12-31-2000", false},
		{"31-12-99", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = IsDate()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isDate(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

// format: HH:MM:SS
func Test_isTime(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"five fifteen", false},
		{"", true},
		{"01:01:01", true},
		{"1:01:01", true},
		{"1:1:01", false},
		{"01:01:010", false},
		{"01 01 01", false},
		{"1:01:1", false},
		{"1:1:1", false},
		{"12:12", false},
		{"12:12:04PM", false},
		{"12:12PM", false},
		{"12:12.12", false},
		{"12.12.12", false},
		{"12-12-12", false},
		{"12+12+12", false},
		{"12_12_12", false},
		{"12:12+12", false},
		{"12/12/12", false},
		{"12\\12\\12", false},
		{"12:12:12", true},
		{"12:12:12 -0700", false},
		{"12:12:12 +0400", false},
		{"00:00:00", true},
		{"23:59:59", true},
		{"23:59:60", false},
		{"12:60:12", false},
		{"24:12:12", false},
		{"83:75:60", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = IsTime()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isTime(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

// DD-MM-YYYY HH:MM:SS
func Test_isDateTime(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"Thursday May 02 2015", false},
		{"", true},
		{"01-01-2001", false},
		{"01:01:01", false},
		{"001-01-2001 01:01:01", false},
		{"01-01-2001 01:01:010", false},
		{"01-01-2001 01:01:01", true},
		{"01-01-2001 1:01:01", true},
		{"01-01-2001 1:1:01", false},
		{"01-01-2001 1:01:1", false},
		{"01-01-2001 1:1:1", false},
		{"01-01-2001 12:12", false},
		{"01-01-2001 12:12:04PM", false},
		{"01-01-2001 12:12PM", false},
		{"01-01-2001 12:12.12", false},
		{"01-01-2001 12.12.12", false},
		{"01-01-2001 12-12-12", false},
		{"01-01-2001 12+12+12", false},
		{"01-01-2001 12_12_12", false},
		{"01-01-2001 12:12+12", false},
		{"01-01-2001 12/12/12", false},
		{"01-01-2001 12\\12\\12", false},
		{"01-01-2001 12:12:12", true},
		{"01-01-2001 12:12:12 -0700", false},
		{"01-01-2001 12:12:12 +0400", false},
		{"01-01-2001 00:00:00", true},
		{"01-01-2001 23:59:59", true},
		{"01-01-2001 23:59:60", false},
		{"01-01-2001 12:60:12", false},
		{"01-01-2001 24:12:12", false},
		{"01-01-2001 83:75:60", false},
		{"29-02-2000 01:01:01", true},
		{"29-02-1999 01:01:01", false},
		{"14-06-2000 01:01:01", true},
		{"14_06_2000 01:01:01", false},
		{"14.06.2000 01:01:01", false},
		{"14/06/2000 01:01:01", false},
		{"14\\06\\2000 01:01:01", false},
		{"32-12-2000 01:01:01", false},
		{"24-13-2000 01:01:01", false},
		{"Monday 22-05-2017 01:01:01", false},
		{"Mon 22-05-2017 01:01:01", false},
		{"Monday 22-May-2017 01:01:01", false},
		{"2000-12-31 01:01:01", false},
		{"2000-31-12 01:01:01", false},
		{"2000/31/12 01:01:01", false},
		{"2000/12/31 01:01:01", false},
		{"99-31-12 01:01:01", false},
		{"24/04/2000 01:01:01", false},
		{"24 04 2000 01:01:01", false},
		{"24 April 2000 01:01:01", false},
		{"24, 04, 2000 01:01:01", false},
		{"24,04,2000 01:01:01", false},
		{"24+04+2000 01:01:01", false},
		{"12-31-2000 01:01:01", false},
		{"31-12-99 01:01:01", false},
	}

	for _, l := range list {
		valid := false
		var rule Rule = IsDateTime()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isDateTime(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isUUID(t *testing.T) {
	var list = []struct {
		field       string
		version     uint32
		expectation bool
	}{
		{"", 0, true}, // unknown version
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", 0, false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3xxx", 0, false},
		{"a987fbc94bed3078cf079141ba07c9f3", 0, false},
		{"934859", 0, false},
		{"987fbc9-4bed-3078-cf07a-9141ba07c9f3", 0, false},
		{"aaaaaaaa-1111-1111-aaag-111111111111", 0, false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", 0, true},

		{"", 3, true}, // v3
		{"412452646", 3, false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", 3, false},
		{"a987fbc9-4bed-4078-8f07-9141ba07c9f3", 3, false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", 3, true},

		{"", 4, true}, // v4
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", 4, false},
		{"a987fbc9-4bed-5078-af07-9141ba07c9f3", 4, false},
		{"934859", 4, false},
		{"57b73598-8764-4ad0-a76a-679bb6640eb1", 4, true},
		{"625e63f3-58f5-40b7-83a1-a72ad31acffb", 4, true},

		{"", 5, true}, // v5
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", 5, false},
		{"9c858901-8a57-4791-81fe-4c455b099bc9", 5, false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", 5, false},
		{"987fbc97-4bed-5078-af07-9141ba07c9f3", 5, true},
		{"987fbc97-4bed-5078-9f07-9141ba07c9f3", 5, true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = IsUUID(l.version)

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("isUUID(%d, '%s'): Valid[%t]. Expected: %t", l.version, l.field, valid, l.expectation)
		}
	}
}

func Test_isbn(t *testing.T) { // no version
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"-45.112023", false},
		{"foo", false},
		{"3423214121", false},
		{"978-3836221191", true},
		{"3-423-21412-1", false},
		{"3 423 21412 1", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = ISBN()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("ISBN(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isbn10(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"-45.112023", false},
		{"foo", false},
		{"3423214121", false},
		{"978-3836221191", false},
		{"3-423-21412-1", false},
		{"3 423 21412 1", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = ISBN10()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("ISBN10(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_isbn13(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"true", false},
		{"10", false},
		{"-5.20301", false},
		{"a", false},
		{"", true},
		{"-45.112023", false},
		{"foo", false},
		{"3-8362-2119-5", false},
		{"01234567890ab", false},
		{"978 3 8362 2119 0", false},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}

	for _, l := range list {
		valid := false
		var rule Rule = ISBN13()

		if e, _ := rule.Validate([]string{l.field}, make(map[string]string)); e == nil {
			valid = true
		}

		if l.expectation != valid {
			t.Errorf("ISBN13(%s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}
