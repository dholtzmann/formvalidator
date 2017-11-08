package formvalidator

import (
	"testing"
)

func Test_appendError(t *testing.T) {
	var a []FormError
	b := &FormError{"Hello", []interface{}{1, 2}}
	c := &FormError{"Goodbye", []interface{}{3, 4}}
	a = appendError(a, b, c)
	if len(a) != 2 {
		t.Errorf("appendError(): Should return a slice with length of 2!")
	}
}

func Test_inSlice(t *testing.T) {
	var list = []struct {
		field       string
		expectation bool
	}{
		{"abc", true},
		{"123", true},
		{"Abc", false},
		{"@bc", false},
		{"ãbc", false},
		{"314", false},
		{"-3.14", false},
	}

	testList := []string{"abc", "def", "ghi", "jkl", "123", "3.14"}
	for _, l := range list {
		valid := inSlice(testList, l.field)
		if l.expectation != valid {
			t.Errorf("inSlice(list, %s): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_hasDuplicates(t *testing.T) {
	var list = []struct {
		field       []string
		expectation bool
	}{
		{[]string{"abc"}, false},
		{[]string{"abc", "Abc"}, false},
		{[]string{"abc", "@bc"}, false},
		{[]string{"abc", "ãbc"}, false},
		{[]string{"abc", "def", "ghi", "jkl"}, false},
		{[]string{"123", "1.23"}, false},
		{[]string{"123", "-123"}, false},
		{[]string{"abc", "abc"}, true},
	}

	for _, l := range list {
		valid := hasDuplicates(l.field)
		if l.expectation != valid {
			t.Errorf("hasDuplicates(%v): Valid[%t]. Expected: %t", l.field, valid, l.expectation)
		}
	}
}

func Test_sliceDiff(t *testing.T) {
	var list = []struct {
		s1          []string
		s2          []string
		expectation []string
	}{
		{[]string{"abc"}, []string{"abc"}, []string{}},
		{[]string{"abc"}, []string{"ãbc"}, []string{"abc"}},
		{[]string{"Abc"}, []string{"abc"}, []string{"Abc"}},
		{[]string{"ãbc"}, []string{"abc"}, []string{"ãbc"}},
		{[]string{"abc", "def", "ghi"}, []string{"abc"}, []string{"def", "ghi"}},
		{[]string{"abc"}, []string{"abc", "def", "ghi"}, []string{}}, // order matters for sliceDiff()
	}

	for _, l := range list {
		result := sliceDiff(l.s1, l.s2)
		if len(l.expectation) != len(result) {
			t.Errorf("sliceDiff(%v, %v): Result: %v, Expected: %v", l.s1, l.s2, result, l.expectation)
		}
	}
}

func Test_mustFormatCSVBytes(t *testing.T) {
	list1 := []byte("abc,def,ghi,jkl")
	list2 := []string{"abc", "def", "ghi", "jkl"}

	result := mustFormatCSVBytes(list1)

	if len(result) != len(list2) {
		t.Errorf("mustFormatCSVBytes(): Result is different from expectation!")
	}
}

func Test_getFirstKey(t *testing.T) {
	list1 := []string{"abc", "def", "ghi"}
	list2 := []string{"123", "456", "789"}

	if getFirstKey(list1) != "abc" {
		t.Errorf("getFirstKey(): Did not return the first value of the provided slice! Return[%s]", getFirstKey(list1))
	}

	if getFirstKey(list2) != "123" {
		t.Errorf("getFirstKey(): Did not return the first value of the provided slice! Return[%s]", getFirstKey(list1))
	}
}
