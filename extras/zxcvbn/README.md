zxcvbn
========

Small extra package for formvalidator, implements formvalidator.Rule.
Checks if a password is at least as strong as the required minimum.

Uses this package to test password strength: [zxcvbn-go](https://github.com/nbutton23/zxcvbn-go)

## Installation

Is installed with formvalidator.

## New Validation rule
IsStrongPassword(minScore uint32, email string)

## Example

```go
	rules := map[string][]fv.Rule{
		"password": fv.RuleChain(IsStrongPassword(2, form.Get("email"))),
	}
```
