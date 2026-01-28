package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type formErrors map[string][]string

var EmailRx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (e formErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e formErrors) Get(field string) string {
	fields, ok := e[field]
	if !ok {
		return ""
	}
	return fields[0]
}

type Form struct {
	url.Values
	formErrors formErrors
}

func NewForm(form url.Values) *Form {
	return &Form{
		Values:     form,
		formErrors: formErrors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) *Form {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			label := []rune(field)
			label[0] = unicode.ToUpper(label[0])
			f.formErrors.Add(field, fmt.Sprintf("%s is required", string(label)))
		}
	}
	return f
}

func (f *Form) Valid() bool {
	return len(f.formErrors) == 0
}

func (f *Form) MinLength(field string, n int) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if utf8.RuneCountInString(value) < n {
		f.formErrors.Add(field, fmt.Sprintf("This field is too short (minimum of %d characters)", n))
	}

	return f
}

func (f *Form) MaxLength(field string, n int) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if utf8.RuneCountInString(value) > n {
		f.formErrors.Add(field, fmt.Sprintf("This field is too long (maximum of %d characters)", n))
	}

	return f
}

func (f *Form) Matches(field string, pattern *regexp.Regexp) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if !pattern.MatchString(value) {
		f.formErrors.Add(field, "This field is invalid")
	}

	return f
}
