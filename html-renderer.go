package formvalidator

import (
	"fmt"
	"errors"
	"html/template"
	"log"
)

/*
Escape and render HTML for the following types: (radio buttons, checkboxes, select-single, and select-multiple)

<input type="radio" value="1" checked="checked">
<input type="checkbox" value="1" checked="checked">
<select><option value="1" selected="selected">1</option></select>
<select multiple="multiple"><option value="1" selected="selected">1</option><option value="2" selected="selected">2</option></select>
*/

/*
	Options are HTML keys
	Labels are Human-readable labels, must have same length as options
	Value/Values are the form field value(s) "<input value='X'>"
	Name is the name of the form field "<input name='X'>"
	CSS is for CSS classes "<input class='center red big-box'>"
	Typ is for the type of form field 'radio','checkbox', or 'select'

	Important:
	Options, Css, Name, and Labels are not escaped for HTML, they are vulnerable to HTML injection/XSS, the reason is that they are supplied by the programmer, not user input!
*/
type singleInput struct {
	name    string
	typ     string
	css     string
	options []string
	labels  []string
	value   string
}

type multiInput struct {
	name    string
	typ     string
	css     string
	options []string
	labels  []string
	values  []string
}

// Errors
var ErrLabelOptionMismatch = errors.New(`Variable 'options []string' must have equal length as 'labels []string'.`)
var ErrLabelDuplicate = errors.New(`Variable 'labels []string has duplicates.'`)
var ErrOptionDuplicate = errors.New(`Variable 'options []string' has duplicates.`)
var ErrBadType = errors.New(`Variable 'typ string' is invalid. Allowed: singleInput['radio','select'] multiInput['checkbox','select']`)
var ErrValueOptionMismatch = errors.New(`multiInput.getHTML(): Variable 'values []string' should have length less than or equal to 'options []string'.`)
var ErrValueDuplicate = errors.New(`multiInput.getHTML(): Variable 'values []string' has duplicates.`)

func NewSingleInput(name, typ, cssClasses string, options, labels []string, value string) *singleInput {
	// add a leading space and class attribute if not blank
	if len(cssClasses) > 0 {
		cssClasses = fmt.Sprintf(` class="%s"`, cssClasses)
	}

	return &singleInput{name, typ, cssClasses, options, labels, value}
}

func NewMultiInput(name, typ, cssClasses string, options, labels, values []string) *multiInput {
	// add a leading space and class attribute if not blank
	if len(cssClasses) > 0 {
		cssClasses = fmt.Sprintf(` class="%s"`, cssClasses)
	}

	return &multiInput{name, typ, cssClasses, options, labels, values}
}

func (i *singleInput) GetHTML() (template.HTML, error) {

	// check for errors
	if len(i.options) != len(i.labels) {
		return template.HTML(""), ErrLabelOptionMismatch
	}
	if hasDuplicates(i.labels) {
		return template.HTML(""), ErrLabelDuplicate
	}
	if hasDuplicates(i.options) {
		return template.HTML(""), ErrOptionDuplicate
	}

	// render it
	switch i.typ {
	case "radio":
		return i.getRadioHTML()

	case "select":
		return i.getSelectHTML()

	default:
		return template.HTML(""), ErrBadType
	}
}

func (i *singleInput) getRadioHTML() (template.HTML, error) {
	var html string

	for n, key := range i.options {
		// it matches, add the "checked" bit
		if i.value == key {
			html += fmt.Sprintf(`<label%s><input type="radio" name="%s" value="%s" checked="checked"> %s</label>`, i.css, i.name, key, i.labels[n])
		} else {
			html += fmt.Sprintf(`<label%s><input type="radio" name="%s" value="%s"> %s</label>`, i.css, i.name, key, i.labels[n])
		}
	}

	return template.HTML(html), nil
}

func (i *singleInput) getSelectHTML() (template.HTML, error) {
	var html, middle string

	for n, key := range i.options {
		// it matches, add the "selected" bit
		if i.value == key {
			middle += fmt.Sprintf(`<option value="%s" selected="selected">%s</option>`, key, i.labels[n])
		} else {
			middle += fmt.Sprintf(`<option value="%s">%s</option>`, key, i.labels[n])
		}
	}

	html = fmt.Sprintf(`<select name="%s"%s>%s</select>`, i.name, i.css, middle)
	return template.HTML(html), nil
}

func (i *multiInput) GetHTML() (template.HTML, error) {

	// check for errors
	if len(i.options) != len(i.labels) {
		return template.HTML(""), ErrLabelOptionMismatch
	}
	if hasDuplicates(i.labels) {
		return template.HTML(""), ErrLabelDuplicate
	}
	if hasDuplicates(i.options) {
		return template.HTML(""), ErrOptionDuplicate
	}
	// i.values could be set by merging raw form data to a form struct
	//	Raw form data from a user cannot be trusted, (HTML injection), but it is not a programmer error so it should not return an error as that would cause a panic with template.Must()
	if len(i.values) > len(i.options) {
		log.Printf(ErrValueOptionMismatch.Error())
		i.values = []string{} // on error set to an empty slice
	}
	if hasDuplicates(i.values) {
		log.Printf(ErrValueDuplicate.Error())
		i.values = []string{} // on error set to an empty slice
	}

	// render it
	switch i.typ {
	case "checkbox":
		return i.getCheckboxHTML()

	case "select":
		return i.getSelectHTML()

	default:
		return template.HTML(""), ErrBadType
	}
}

func (i *multiInput) getCheckboxHTML() (template.HTML, error) {
	var html string

	for n, key := range i.options {
		// match found, add the "checked" bit
		if inSlice(i.values, key) {
			html += fmt.Sprintf(`<label%s><input type="checkbox" name="%s" value="%s" checked="checked"> %s</label>`, i.css, i.name, key, i.labels[n])
		} else {
			html += fmt.Sprintf(`<label%s><input type="checkbox" name="%s" value="%s"> %s</label>`, i.css, i.name, key, i.labels[n])
		}
	}

	return template.HTML(html), nil
}

func (i *multiInput) getSelectHTML() (template.HTML, error) {
	var html, middle string

	for n, key := range i.options {
		// match found, add the "checked" bit
		if inSlice(i.values, key) {
			middle += fmt.Sprintf(`<option value="%s" selected="selected">%s</option>`, key, i.labels[n])
		} else {
			middle += fmt.Sprintf(`<option value="%s">%s</option>`, key, i.labels[n])
		}
	}

	html = fmt.Sprintf(`<select name="%s"%s multiple="multiple">%s</select>`, i.name, i.css, middle)
	return template.HTML(html), nil
}
