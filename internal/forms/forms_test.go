package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}
  form := New(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
  form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ := http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
  form := New(postedData)
	r := httptest.NewRequest("POST", "/whatever", nil)

	has := form.Has("whatever", r)

	if has {
		t.Error("shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("key", "values")

	form = New(postedData)

	has = form.Has("key", r)

	if !has {
		t.Error("does not show has field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
  form := New(postedData)
	r := httptest.NewRequest("POST", "/whatever", nil)

	form.MinLength("x", 10, r)
	if form.Valid() {
		t.Error("valid for not exsisting attribute")
	}

	postedData = url.Values{}
	postedData.Add("key", "values")
	form = New(postedData)
	form.MinLength("key", 1, r)
	if !form.Valid() {
		t.Error("valid for valid length")
	}
}


func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
  form := New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("valid for not exsisting attribute")
	}

	postedData = url.Values{}
	postedData.Add("email", "not_valid")
	form = New(postedData)
	form.IsEmail("email")

	if form.Valid() {
		t.Error("valid for not valid email")
	}

	postedData.Add("valid_email", "test@gmail.com")
	form = New(postedData)
	form.IsEmail("valid_email")

	if !form.Valid() {
		t.Error("not valid for valid email")
	}
}
