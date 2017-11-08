package formvalidator

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

/*
	Custom error type

	Str can be accessed to do translation, it has the raw data so Sprintf flags are included (%d,%s,etc.)

	Example:

	var e = FormError{"Integer must be between %d and %d", []interface{}{5, 10}}
	i18n(e.Str, e.Data) -> "Escriba Usted un valor entre 5 y 10"
*/
type FormError struct {
	Str  string
	Data []interface{}
}

func (e *FormError) Error() string {
	if len(e.Data) > 0 { // if 'Data' is not empty format the string
		return fmt.Sprintf(e.Str, e.Data...)
	}
	return e.Str
}

// appendError(err, errors...) do not forget the dots for slices
func appendError(err []FormError, arg ...*FormError) []FormError {
	for _, e := range arg {
		err = append(err, *e)
	}
	return err
}

// Is some value in the slice? (string type only)
func inSlice(slice []string, val string) bool {
	for _, j := range slice {
		if j == val {
			return true
		}
	}
	return false
}

// check for duplicate entries in a slice
func hasDuplicates(slice []string) bool {
	encounteredItems := make(map[string]bool)
	for _, val := range slice {

		// new entry, add it to the found list
		if _, found := encounteredItems[val]; !found {
			encounteredItems[val] = true
		} else {
			// already has item set in map, duplicate found
			return true
		}
	}
	// no duplicates found
	return false
}

// Compare two slices and return the difference (string types only for slices)
func sliceDiff(slice1, slice2 []string) (difference []string) {
	for _, j := range slice1 {
		if !inSlice(slice2, j) {
			difference = append(difference, j)
		}
	}
	return
}

// some validator rules load data from a csv file
func mustFormatCSVBytes(b []byte) []string {
	buf := bytes.NewBuffer(b)
	r := csv.NewReader(buf)
	var record [][]string
	var err error
	if record, err = r.ReadAll(); err != nil {
		panic(err.Error())
	}

	var list []string
	if len(record) == 1 { // the file should have only one line, the multi-dimensional slice is for rows and columns, see docs on "encoding/csv"
		list = record[0]
	}

	return list
}

// get the first key from a slice, []string
func getFirstKey(values []string) string {
	if len(values) > 0 {
		return values[0]
	}

	return ""
}
