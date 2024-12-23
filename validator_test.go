package structvalidator

import (
	"log"
	"testing"
)

type Test1 struct {
	FirstName     string `validation:"req lenmin:5 lenmax:25"`
	LastName      string `validation:"req lenmin:2 lenmax:50"`
	Age           int    `validation:"req valmin:18 valmax:150"`
	Price         int    `validation:"req valmin:0 valmax:9999"`
	PostCode      string `validation:"req" validation_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `validation:"req email"`
	BelowZero     int    `validation:"valmin:-6 valmax:-2"`
	DiscountPrice int    `validation:"valmin:0 valmax:8000"`
	Country       string `validation_regexp:"^[A-Z][A-Z]$"`
	County        string `validation:"lenmax:40 uitextarea"`
}

type Test2 struct {
	FirstName     string `mytag:"req lenmin:5 lenmax:25"`
	LastName      string `mytag:"req lenmin:2 lenmax:50"`
	Age           int    `mytag:"req valmin:18 valmax:150"`
	Price         int    `mytag:"req valmin:0 valmax:9999"`
	PostCode      string `mytag:"req" mytag_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `mytag:"req email"`
	BelowZero     int    `mytag:"valmin:-6 valmax:-2"`
	DiscountPrice int    `mytag:"valmin:0 valmax:8000"`
	Country       string `mytag_regexp:"^[A-Z][A-Z]$"`
	County        string `mytag:"lenmax:40"`
}

type Test3 struct {
	ZeroMin  int `mytag:"valmin:0 valmax:5"`
	ZeroMax  int `mytag:"valmax:0"`
	ZeroBoth int `mytag:"valmin:0 valmax:0"`
	NotZero  int `mytag:"valmin:4 valmax:6"`
	OnlyMin  int `mytag:"valmin:3"`
	OnlyMax  int `mytag:"valmax:7"`
}

type Test4 struct {
	PrimaryEmail string ``
}

func TestWithDefaultValues(t *testing.T) {
	s := Test1{}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"FirstName": FailEmpty,
		"LastName":  FailEmpty,
		"Age":       FailValMin,
		"PostCode":  FailEmpty,
		"Email":     FailEmpty,
		"Country":   FailRegexp,
		"BelowZero": FailValMax,
	}
	opts := &ValidationOptions{}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithInvalidValues(t *testing.T) {
	s := Test1{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
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
	opts := &ValidationOptions{}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithValidValues(t *testing.T) {
	s := Test1{
		FirstName:     "Johnny",
		LastName:      "Smith",
		Age:           35,
		Price:         0,
		PostCode:      "43-155",
		Email:         "john@example.com",
		BelowZero:     -4,
		DiscountPrice: 8000,
		Country:       "GB",
		County:        "Enfield",
	}
	expectedBool := true
	expectedFailedFields := map[string]int{}
	opts := &ValidationOptions{}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithInvalidValuesAndFieldRestriction(t *testing.T) {
	s := Test1{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
	}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"FirstName": FailLenMax,
		"LastName":  FailLenMin,
	}
	opts := &ValidationOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"LastName":  true,
		},
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithInvalidValuesAndFieldRestrictionAndOverwrittenFieldTags(t *testing.T) {
	s := Test1{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
	}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"LastName": FailLenMin,
	}
	opts := &ValidationOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"LastName":  true,
		},
		OverwriteFieldTags: map[string]map[string]string{
			"FirstName": map[string]string{
				"validation": "req lenmin:4 lenmax:100",
			},
		},
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithInvalidValuesAndOverwrittenTagName(t *testing.T) {
	s := Test2{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
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
		OverwriteTagName: "mytag",
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestValMinMaxWithDefault(t *testing.T) {
	s := Test3{}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"NotZero": FailValMin,
		"OnlyMin": FailValMin,
	}
	opts := &ValidationOptions{
		OverwriteTagName: "mytag",
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestValMinMaxWithValid(t *testing.T) {
	s := Test3{
		NotZero: 4,
		OnlyMin: 3,
		OnlyMax: 7,
	}
	expectedBool := true
	expectedFailedFields := map[string]int{}
	opts := &ValidationOptions{
		OverwriteTagName: "mytag",
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestValMinMaxWithInvalid(t *testing.T) {
	s := Test3{
		ZeroMin:  -4,
		ZeroMax:  -6,
		ZeroBoth: -6,
		NotZero:  2,
		OnlyMin:  -5,
		OnlyMax:  -6,
	}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"ZeroMin":  FailValMin,
		"ZeroBoth": FailValMin,
		"NotZero":  FailValMin,
		"OnlyMin":  FailValMin,
	}
	opts := &ValidationOptions{
		OverwriteTagName: "mytag",
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithInvalidValuesWithSuffixValidation(t *testing.T) {
	s := Test4{
		PrimaryEmail: "invalidemail",
	}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"PrimaryEmail": FailEmail,
	}
	opts := &ValidationOptions{
		ValidateWhenSuffix: true,
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithInvalidValuesWithoutSuffixValidation(t *testing.T) {
	s := Test4{
		PrimaryEmail: "invalidemail",
	}
	expectedBool := true
	expectedFailedFields := map[string]int{}
	opts := &ValidationOptions{
		ValidateWhenSuffix: false,
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func TestWithOverwrittenValues(t *testing.T) {
	s := Test1{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
	}
	expectedBool := false
	expectedFailedFields := map[string]int{
		"Age": FailValMax,
	}
	opts := &ValidationOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"LastName":  true,
			"Age":       true,
		},
		OverwriteFieldValues: map[string]interface{}{
			"FirstName": "123456",
			"LastName":  "123",
			"Age":       300,
		},
	}
	compare(&s, expectedBool, expectedFailedFields, opts, t)
}

func compare(s interface{}, expectedBool bool, expectedFailedFields map[string]int, options *ValidationOptions, t *testing.T) {
	valid, failedFields := Validate(s, options)
	if valid != expectedBool {
		t.Fatalf("Validate returned invalid boolean value")
	}
	compareFailedFields(failedFields, expectedFailedFields, t)
}

func compareFailedFields(failedFields map[string]int, expectedFailedFields map[string]int, t *testing.T) {
	if len(failedFields) != len(expectedFailedFields) {
		for k, v := range failedFields {
			log.Printf("%s %d", k, v)
		}
		t.Fatalf("Validate returned invalid number of failed fields %d where it should be %d", len(failedFields), len(expectedFailedFields))
	}
	for k, v := range expectedFailedFields {
		if failedFields[k] != v {
			t.Fatalf("Validate returned invalid failure flag of %d where it should be %d for %s", failedFields[k], v, k)
		}
	}
}
