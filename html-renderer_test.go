package formvalidator

import (
	"strings"
	"testing"
)

// TODO add more tests (error types from htmlrenderer.go)

func Test_singleInputRadioHTMLInjection(t *testing.T) {
	injection := `javascript:"alert('Hello');"`
	i := NewSingleInput("email", "radio", "", []string{"gmail.com", "hotmail.com", "yahoo.com"}, []string{"G-mail", "Hotmail", "Yahoo!"}, injection)
	html, err := i.GetHTML()

	if err != nil {
		t.Errorf("SingleInput Error: %s\n", err.Error())
	}

	if strings.Contains(string(html), injection) == true {
		t.Errorf("SingleInput [Type:radio] value is vulnerable to html injection!")
	}
}

func Test_singleInputSelectHTMLInjection(t *testing.T) {
	injection := `javascript:"alert('Hello');"`
	i := NewSingleInput("email", "select", "", []string{"gmail.com", "hotmail.com", "yahoo.com"}, []string{"G-mail", "Hotmail", "Yahoo!"}, injection)
	html, err := i.GetHTML()

	if err != nil {
		t.Errorf("SingleInput Error: %s\n", err.Error())
	}

	if strings.Contains(string(html), injection) == true {
		t.Errorf("SingleInput [Type:select] value is vulnerable to html injection!")
	}
}

func Test_multiInputCheckboxHTMLInjection(t *testing.T) {
	injection := []string{`javascript:"alert('Hello');"`, `gmail.com`}
	i := NewMultiInput("email", "checkbox", "", []string{"gmail.com", "hotmail.com", "yahoo.com"}, []string{"G-mail", "Hotmail", "Yahoo!"}, injection)
	html, err := i.GetHTML()

	if err != nil {
		t.Errorf("MultiInput Error: %s\n", err.Error())
	}

	if strings.Contains(string(html), injection[0]) == true {
		t.Errorf("MultiInput [Type:checkbox] value is vulnerable to html injection!")
	}
}

func Test_multiInputSelectHTMLInjection(t *testing.T) {
	injection := []string{`javascript:"alert('Hello');"`, `gmail.com`}
	i := NewMultiInput("email", "select", "", []string{"gmail.com", "hotmail.com", "yahoo.com"}, []string{"G-mail", "Hotmail", "Yahoo!"}, injection)
	html, err := i.GetHTML()

	if err != nil {
		t.Errorf("MultiInput Error: %s\n", err.Error())
	}

	if strings.Contains(string(html), injection[0]) == true {
		t.Errorf("MultiInput [Type:select] value is vulnerable to html injection!")
	}
}
