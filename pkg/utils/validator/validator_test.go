package validator

import (
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func TestNewValidator_UUIDValidation(t *testing.T) {
	v := NewValidator()

	type S struct {
		ID string `validate:"uuid"`
	}
	var s S

	s.ID = uuid.New().String()
	err := v.Struct(s)
	if err != nil {
		t.Fatalf("expected no error for valid uuid, got: %v", err)
	}

	s.ID = "not-a-uuid"
	err = v.Struct(s)
	if err == nil {
		t.Fatalf("expected error for invalid uuid, got nil")
	}
}

func TestValidatorErrors(t *testing.T) {
	if msg := ValidatorErrors(nil); msg != "" {
		t.Errorf("expected empty string for nil error, got %q", msg)
	}

	err := validator.ValidationErrors{}
	if msg := ValidatorErrors(err); msg != "" {
		t.Errorf("expected empty string for empty ValidationErrors, got %q", msg)
	}

	if msg := ValidatorErrors(assertAnError{}); msg != "" {
		t.Errorf("expected empty string for non-validation error, got %q", msg)
	}

	errs := validator.ValidationErrors{fakeErr{"Foo", "required"}, fakeErr{"Bar", "uuid"}}
	msg := ValidatorErrors(errs)
	if msg != "Foo: required\nBar: uuid" {
		t.Errorf("unexpected ValidatorErrors output: %q", msg)
	}
}

type fakeErr struct{ field, tag string }

func (f fakeErr) Field() string                    { return f.field }
func (f fakeErr) ActualTag() string                { return f.tag }
func (f fakeErr) Namespace() string                { return "" }
func (f fakeErr) StructNamespace() string          { return "" }
func (f fakeErr) StructField() string              { return "" }
func (f fakeErr) Tag() string                      { return "" }
func (f fakeErr) Kind() reflect.Kind               { return reflect.String }
func (f fakeErr) Type() reflect.Type               { return reflect.TypeOf("") }
func (f fakeErr) Value() interface{}               { return "" }
func (f fakeErr) Param() string                    { return "" }
func (f fakeErr) Error() string                    { return f.field + ": " + f.tag }
func (f fakeErr) Translate(_ ut.Translator) string { return f.Error() }

type assertAnError struct{}

func (assertAnError) Error() string { return "assert error" }
