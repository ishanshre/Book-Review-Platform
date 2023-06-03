package forms

import (
	"fmt"
	"net/url"
	"strings"
)

// From type holds the data from forms and errors
type Form struct {
	url.Values
	Errors errors
}

// New initializes an empty Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks all the required fields passed to it
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength checks the minimum length of characters in the field
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long.", length))
		return false
	}
	return true
}

// Has checks if the form field is empty or not
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be black")
		return false
	}
	return true
}

// Valid reutrns true if there are no errors else returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
