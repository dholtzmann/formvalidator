package zxcvbn

import (
	"errors"

	zgo "github.com/nbutton23/zxcvbn-go"
	"github.com/nbutton23/zxcvbn-go/matching"
	"github.com/nbutton23/zxcvbn-go/scoring"
)

// for validating an HTML form field, implements an interface from another package
type isStrongPassword struct {
	minScore uint32
	userInputs []string
}

/*
	Score has a range of 0-4, see (https://github.com/nbutton23/zxcvbn-go) for details.
	E-mail address can be added to verify that it is not used as the password.	
*/
func IsStrongPassword(minScore uint32, email string) (isStrongPassword) {
	if minScore > 4 {
		minScore = 4
	}

	if len(email) > 250 { // truncate e-mail string if too long
		email = email[0:250]
	}
	return isStrongPassword{minScore, []string{email}}
}

/*
	Check if the submitted password is at least as strong as the required minimum, uses zxcvbn-go to test strength
*/
func (i isStrongPassword) Validate(fields []string, errorMessages map[string]string) (error, []interface{}) {

	var errMsg = errors.New(errorMessages["weak_password"])
	var field = ""
	var result scoring.MinEntropyMatch

	if len(fields) > 0 {
		field = fields[0]
	}


	/*
		zxcvbn-go has a bug, if a long password (25+ characters) with "L33t" characters is used, like an MD5 sum, it allocates a lot of memory (in the gigabytes!) and takes too long (hours!).

		See the Issue titled: "zxcvbn pathological case" on https://github.com/nbutton23/zxcvbn-go

		For a workaround, check for a relatively long password and use the matching.FilterL33tMatcher flag.
	*/
	if len(field) < 17 {
		result = zgo.PasswordStrength(field, i.userInputs)
	} else {
		result = zgo.PasswordStrength(field, i.userInputs, matching.FilterL33tMatcher)
	}

	if result.Score >= int(i.minScore) {
		return nil, nil
	}
 
	return errMsg, nil
}
