package structvalidator

import (
	"reflect"
	"testing"
)

// In the struct-db-postgres, struct-validator is used to validate map of values against reflect.Value which is actually
// a pointer of a pointer to struct.
// TODO: This should be revisited at some point.
type Wrapper struct {
	DoesntMatter    string
	UseMeToValidate []*Test1
}

func TestWithInvalidValuesOnReflectValue(t *testing.T) {
	o := &Wrapper{}
	v := reflect.ValueOf(o)
	i := reflect.Indirect(v)
	s := i.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		k := f.Type.Kind()

		// Only field which are slices of pointers to struct instances
		if k != reflect.Slice || f.Type.Elem().Kind() != reflect.Ptr || f.Type.Elem().Elem().Kind() != reflect.Struct {
			continue
		}

		expectedBool := false
		expectedFailedFields := map[string]int{
			"FirstName":     FailLenMax,
			"LastName":      FailLenMin,
			"Age":           FailValMin,
			"PostCode":      FailRegexp,
			"Email":         FailEmail,
			"BelowZero":     FailValMax,
			"DiscountPrice": FailValMax,
			"Country":       FailRegexp,
		}
		opts := &ValidationOptions{
			OverwriteFieldValues: map[string]interface{}{
				"FirstName":     "123456789012345678901234567890",
				"LastName":      "b",
				"Age":           15,
				"Price":         0,
				"PostCode":      "AA123",
				"Email":         "invalidEmail",
				"BelowZero":     8,
				"DiscountPrice": 9999,
				"Country":       "Tokelau",
				"County":        "",
			},
		}
		compare(reflect.New(f.Type.Elem()), expectedBool, expectedFailedFields, opts, t)
	}
}
