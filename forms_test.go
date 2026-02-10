package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewForm(t *testing.T) {
	assert := assert.New(t) // Initialize once

	values := url.Values{}
	values.Add("email", "test@mail.com")

	form := NewForm(values)
	assert.NotNil(form)
	assert.Equal("test@mail.com", values.Get("email"))
	assert.NotNil(form.Errors)
	assert.Len(form.Errors, 0)
}

func TestFormRequired(t *testing.T) {
	values := url.Values{}
	values.Add("email", "test@mail.com")
	values.Add("empty", "  ")

	form := NewForm(values)
	form.Required("email", "password", "empty")

	assert.NotNil(t, form)
	assert.Equal(t, "", form.Errors.Get("email"))
	assert.Contains(t, form.Errors.Get("password"), "Password is required")
	assert.Contains(t, form.Errors.Get("empty"), "Empty is required")
}
