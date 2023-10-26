package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
)

type Form struct {
	url.Values
	Errors errros
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errros(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

func (f *Form) MinLength(field string, len int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) < len {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d)", len))
	}
}

func (f *Form) MatchPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
