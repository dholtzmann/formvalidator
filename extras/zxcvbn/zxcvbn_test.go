package zxcvbn

import (
	"testing"
	"net/url"

	fv "github.com/dholtzmann/formvalidator"
)

func Test_formValidate(t *testing.T) {
	form := url.Values{}
	form.Set("email", "osierjlaskdfjosierjlaksdjfoisajdf@oriseorjsdf.com")
	form.Set("password", "e2c569be17396eca2a2e3c11578123ede2c569be17396eca2a2e3c1$^*:>{<>#$!$~1578123ede2c569be17396eca2a2e3c11578123ede2c569be17396eca2a2e3c11578123ede2c569be17396eca2a2e3c11578123ed")

	// this should validate without a problem
	rules := map[string][]fv.Rule{
		"password": fv.RuleChain(IsStrongPassword(2, form.Get("email"))),
	}

	err, validator := fv.New(rules)
	if err != nil {
		t.Errorf("Error making a new form validator type! %s", err.Error())
	}

	_, errors := validator.Validate(form)

	if len(errors) > 0 {
		for field, val := range errors {
			for _, e := range val {
				t.Errorf("Test_formValidate(): %s: %s", field, e.Error())
			}
		}
	}
}
