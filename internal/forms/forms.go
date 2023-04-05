package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values // READ anonymus interfaces !!! "url.Values": This is an anonymous field of type "url.Values", which is an interface type that represents a collection of key/value pairs
	Errors errors
}


func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}


func (f *Form) Has(field string, r *http.Request) bool {
	x := f.Get(field)

	if x == "" {
		f.Errors.Add(field, "This field can not be blank")
		return false
	}

	return true
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be blank")
		}
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d long", length))
		return false
	}

	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email Adress")
	}
}
