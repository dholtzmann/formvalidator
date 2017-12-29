zxcvbn
========

Small extra package for formvalidator, implements interface formvalidator.Rule

Checks if a password is at least as strong as the required minimum.

Uses [zxcvbn-go](https://github.com/nbutton23/zxcvbn-go) to test password strength.

## Installation

```bash
go get -u github.com/dholtzmann/formvalidator
```

## New Validation rule
- IsStrongPassword(minScore uint32, email string)

## Example

```go
	import(
		fv "github.com/dholtzmann/formvalidator"
		"github.com/dholtzmann/formvalidator/extras/zxcvbn"
	)

	// ...

	func something() {
		rules := map[string][]fv.Rule{
			"password": fv.RuleChain(zxcvbn.IsStrongPassword(2, form.Get("email"))),
		}
	}
```
