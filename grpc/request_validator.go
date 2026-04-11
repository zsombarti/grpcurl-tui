package grpc

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// ValidationError holds a field path and the reason it failed validation.
type ValidationError struct {
	Field  string
	Reason string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("field %q: %s", e.Field, e.Reason)
}

// RequestValidator validates a dynamic proto message against its descriptor.
type RequestValidator struct{}

// NewRequestValidator returns a new RequestValidator.
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// Validate checks required fields on the given protoreflect.Message.
// It returns a slice of ValidationError for every violation found.
func (v *RequestValidator) Validate(msg protoreflect.Message) []ValidationError {
	if msg == nil {
		return []ValidationError{{Field: "(root)", Reason: "message is nil"}}
	}
	var errs []ValidationError
	v.walk(msg, "", &errs)
	return errs
}

// IsValid returns true when Validate produces no errors.
func (v *RequestValidator) IsValid(msg protoreflect.Message) bool {
	return len(v.Validate(msg)) == 0
}

// Summary returns a single joined error string or empty string when valid.
func (v *RequestValidator) Summary(msg protoreflect.Message) error {
	errs := v.Validate(msg)
	if len(errs) == 0 {
		return nil
	}
	parts := make([]string, len(errs))
	for i, e := range errs {
		parts[i] = e.Error()
	}
	return errors.New(strings.Join(parts, "; "))
}

func (v *RequestValidator) walk(msg protoreflect.Message, prefix string, errs *[]ValidationError) {
	msg.Range(func(fd protoreflect.FieldDescriptor, val protoreflect.Value) bool {
		path := fieldPath(prefix, string(fd.Name()))
		if fd.Kind() == protoreflect.MessageKind && !fd.IsList() && !fd.IsMap() {
			v.walk(val.Message(), path, errs)
		}
		return true
	})

	fields := msg.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if fd.Cardinality() == protoreflect.Required && !msg.Has(fd) {
			path := fieldPath(prefix, string(fd.Name()))
			*errs = append(*errs, ValidationError{Field: path, Reason: "required field missing"})
		}
	}
}

func fieldPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}
