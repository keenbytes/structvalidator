package structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// values for invalid field flags
const (
	_          = iota
	FailLenMin = 1 << iota
	FailLenMax
	FailValMin
	FailValMax
	FailEmpty
	FailRegexp
	FailEmail
	FailZero
)

// Optional configuration for validation:
// * RestrictFields defines what struct fields should be validated
// * OverwriteFieldTags can be used to overwrite tags for specific fields
// * OverwriteTagName sets tag used to define validation (default is "validation")
// * ValidateWhenSuffix will validate certain fields based on their name, eg. "PrimaryEmail" field will need to be a valid email
// * OverwriteFieldValues is to use overwrite values for fields, so these values are validated not the ones in struct
type ValidationOptions struct {
	RestrictFields       map[string]bool
	OverwriteFieldTags   map[string]map[string]string
	OverwriteTagName     string
	ValidateWhenSuffix   bool
	OverwriteFieldValues map[string]interface{}
}

// Validate validates fields of a struct.  Currently only fields which are string or int (any) are validated.
// Func returns boolean value that determines whether value is true or false, and a map of fields that failed
// validation.  See Fail* constants for the values.
func Validate(obj interface{}, options *ValidationOptions) (bool, map[string]int) {
	// ValidationOptions is required
	if options == nil {
		panic("ValidationOptions cannot be nil")
	}

	v := reflect.ValueOf(obj)
	i := reflect.Indirect(v)
	s := i.Type()

	// TODO: Fix this to traverse the pointer behind reflect.Value properly.  Current this is made to support
	// struct-db-postgres module that uses this validator.
	if s.String() == "reflect.Value" {
		s = reflect.ValueOf(obj.(reflect.Value).Interface()).Type().Elem().Elem()
	}

	tagName := "validation"
	if options.OverwriteTagName != "" {
		tagName = options.OverwriteTagName
	}

	invalidFields := make(map[string]int, s.NumField())
	valid := true

	for j := 0; j < s.NumField(); j++ {
		field := s.Field(j)
		fieldKind := field.Type.Kind()

		// check if only specified field should be checked
		if len(options.RestrictFields) > 0 && !options.RestrictFields[field.Name] {
			continue
		}

		// validate only ints and string
		if !isInt(fieldKind) && fieldKind != reflect.String {
			continue
		}

		validation := NewValueValidation()

		tagVal, tagRegexpVal := getFieldTagValues(&field, tagName, options.OverwriteFieldTags)
		setValidationFromTags(validation, tagVal, tagRegexpVal)
		if options.ValidateWhenSuffix {
			setValidationFromSuffix(validation, &field)
		}

		// field value can be overwritten in ValidationOptions
		var fieldValue reflect.Value
		overwriteVal, ok := options.OverwriteFieldValues[field.Name]
		if ok {
			fieldValue = reflect.ValueOf(overwriteVal)
		} else {
			fieldValue = v.Elem().FieldByName(field.Name)
		}

		ok, failureFlags := validation.ValidateReflectValue(fieldValue)
		if !ok {
			valid = false
			invalidFields[field.Name] = failureFlags
		}
	}

	return valid, invalidFields
}

func setValidationFromTags(v *ValueValidation, tag string, tagRegexp string) {
	opts := strings.SplitN(tag, " ", -1)
	for _, opt := range opts {
		if opt == "req" {
			v.Flags = v.Flags | Required
		}
		if opt == "email" {
			v.Flags = v.Flags | Email
		}
		for _, valOpt := range []string{"lenmin", "lenmax", "valmin", "valmax", "regexp"} {
			if strings.HasPrefix(opt, valOpt+":") {
				val := strings.Replace(opt, valOpt+":", "", 1)
				if valOpt == "regexp" {
					v.Regexp = regexp.MustCompile(val)
					continue
				}

				i, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				switch valOpt {
				case "lenmin":
					v.LenMin = i
				case "lenmax":
					v.LenMax = i
				case "valmin":
					v.ValMin = int64(i)
					if i == 0 {
						v.Flags = v.Flags | ValMinNotNil
					}
				case "valmax":
					v.ValMax = int64(i)
					if i == 0 {
						v.Flags = v.Flags | ValMaxNotNil
					}
				}
			}
		}
	}

	if tagRegexp != "" {
		v.Regexp = regexp.MustCompile(tagRegexp)
	}
}

func setValidationFromSuffix(v *ValueValidation, field *reflect.StructField) {
	if strings.HasSuffix(field.Name, "Email") {
		v.Flags = v.Flags | Email
	}
	if strings.HasSuffix(field.Name, "Price") && v.ValMin == 0 && v.ValMax == 0 && v.Flags&ValMinNotNil == 0 && v.Flags&ValMaxNotNil == 0 {
		v.ValMin = 0
		v.Flags = v.Flags | ValMinNotNil
	}
}

func isInt(k reflect.Kind) bool {
	if k == reflect.Int64 || k == reflect.Int32 || k == reflect.Int16 || k == reflect.Int8 || k == reflect.Int || k == reflect.Uint64 || k == reflect.Uint32 || k == reflect.Uint16 || k == reflect.Uint8 || k == reflect.Uint {
		return true
	}
	return false
}

func getFieldTagValues(field *reflect.StructField, tagName string, overwriteFieldTags map[string]map[string]string) (tagVal string, tagRegexpVal string) {
	tagVal = field.Tag.Get(tagName)
	tagRegexpVal = field.Tag.Get(tagName + "_regexp")

	overwriteTags, ok := overwriteFieldTags[field.Name]
	if ok {
		overwriteTagVal, ok2 := overwriteTags[tagName]
		if ok2 {
			tagVal = overwriteTagVal
		}
		overwriteTagVal, ok2 = overwriteTags[tagName+"_regexp"]
		if ok2 {
			tagRegexpVal = overwriteTagVal
		}
	}
	return
}
