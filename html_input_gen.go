package structvalidator

import (
	"fmt"
	"html"
	"reflect"
	"strconv"
	"strings"
)

const TypeText = 1
const TypeTextarea = 2
const TypePassword = 3
const TypeEmail = 4

// Optional configuration for validation:
// * RestrictFields defines what struct fields should be generated
// * ExcludeFields defines fields that should be skipped (also from RestrictFields)
// * OverwriteFieldTags can be used to overwrite tags for specific fields
// * OverwriteTagName sets tag used to define validation (default is "validation")
// * ValidateWhenSuffix will validate certain fields based on their name, eg. "PrimaryEmail" field will need to be a valid email
// * OverwriteFieldValues is to use overwrite values for fields, so these values are validated not the ones in struct
// * IDPrefix - if added, an element will contain an 'id' attribute in form of prefix + field name
// * NamePrefix - use this to put a prefix in the 'name' attribute
// * OverwriteValues - fill inputs with the specified values
// * FieldValues - when true then fill inputs with struct instance values
type HTMLOptions struct {
	RestrictFields     map[string]bool
	ExcludeFields      map[string]bool
	OverwriteFieldTags map[string]map[string]string
	OverwriteTagName   string
	ValidateWhenSuffix bool
	IDPrefix           string
	NamePrefix         string
	OverwriteValues    map[string]string
	FieldValues        bool
}

// GenerateHTMLInput takes a struct and generates HTML inputs for each of the fields, eg. <input> or <textarea>
func GenerateHTML(obj interface{}, options *HTMLOptions) map[string]string {
	v := reflect.ValueOf(obj)
	i := reflect.Indirect(v)
	s := i.Type()
	elem := v.Elem()

	tagName := "validation"
	if options != nil && options.OverwriteTagName != "" {
		tagName = options.OverwriteTagName
	}

	fields := map[string]string{}

	for j := 0; j < s.NumField(); j++ {
		field := s.Field(j)
		fieldKind := field.Type.Kind()

		// check if only specified field should be checked
		if options != nil && len(options.RestrictFields) > 0 && !options.RestrictFields[field.Name] {
			continue
		}

		// check if field should not be excluded
		if options != nil && len(options.ExcludeFields) > 0 && options.ExcludeFields[field.Name] {
			continue
		}

		// generate only ints, string and bool
		if !isInt(fieldKind) && !isString(fieldKind) && !isBool(fieldKind) {
			continue
		}

		// value
		value := ""

		if options != nil && options.FieldValues {
			if isBool(fieldKind) && elem.Field(j).Bool() {
				value = "true"
			}
			if isString(fieldKind) {
				value = elem.Field(j).String()
			}
			if isInt(fieldKind) {
				value = fmt.Sprintf("%d", elem.Field(j).Int())
			}
		}

		if options != nil && len(options.OverwriteValues) > 0 && options.OverwriteValues[field.Name] != "" {
			value = options.OverwriteValues[field.Name]
		}

		// 'id' attribute
		fieldIDAttr := ""
		if options.IDPrefix != "" {
			fieldIDAttr = fmt.Sprintf(` id="%s%s"`, options.IDPrefix, field.Name)
		}
		fieldNameAttr := fmt.Sprintf(` name="%s%s"`, options.NamePrefix, field.Name)

		if isBool(fieldKind) {
			fieldChecked := ""
			if value == "true" {
				fieldChecked = " checked"
			}

			fields[field.Name] = fmt.Sprintf(`<input type="checkbox"%s%s%s>`, fieldNameAttr, fieldIDAttr, fieldChecked)
			continue
		}

		// get tag values
		tagVal := field.Tag.Get(tagName)
		tagRegexpVal := field.Tag.Get(tagName + "_regexp")
		if options != nil && len(options.OverwriteFieldTags) > 0 {
			if len(options.OverwriteFieldTags[field.Name]) > 0 {
				if options.OverwriteFieldTags[field.Name][tagName] != "" {
					tagVal = options.OverwriteFieldTags[field.Name][tagName]
				}
				if options.OverwriteFieldTags[field.Name][tagName+"_regexp"] != "" {
					tagRegexpVal = options.OverwriteFieldTags[field.Name][tagName+"_regexp"]
				}
			}
		}

		patternAttr := ""
		if tagRegexpVal != "" {
			patternAttr = fmt.Sprintf(` pattern="%s"`, html.EscapeString(tagRegexpVal))
		}
		validationAttrs, inputType := getHTMLAttributesFromTag(tagVal)

		if options != nil && options.ValidateWhenSuffix {
			if strings.HasSuffix(field.Name, "Email") {
				inputType = TypeEmail
			}
			// Price not supported here yet
		}

		fieldValue := ""
		if value != "" {
			fieldValue = fmt.Sprintf(" value=\"%s\"", html.EscapeString(value))
		}

		if isInt(fieldKind) {
			fields[field.Name] = fmt.Sprintf(`<input type="number"%s%s%s%s/>`, fieldNameAttr, fieldIDAttr, validationAttrs, fieldValue)
			continue
		}

		if isString(fieldKind) {
			if inputType == TypeTextarea {
				fields[field.Name] = fmt.Sprintf(`<textarea%s%s%s%s>%s</textarea>`, fieldNameAttr, fieldIDAttr, validationAttrs, patternAttr, html.EscapeString(value))
				continue
			}
			fieldTypeAttr := ` type="text"`
			if inputType == TypeEmail {
				fieldTypeAttr = ` type="email"`
			}
			if inputType == TypePassword {
				fieldTypeAttr = ` type="password"`
				fieldValue = ""
			}
			fields[field.Name] = fmt.Sprintf(`<input%s%s%s%s%s%s/>`, fieldTypeAttr, fieldNameAttr, fieldIDAttr, validationAttrs, patternAttr, fieldValue)
			continue
		}
	}

	return fields
}

func getHTMLAttributesFromTag(tag string) (string, int) {
	attrs := ""
	inputType := TypeText

	opts := strings.SplitN(tag, " ", -1)
	for _, opt := range opts {
		if opt == "req" {
			attrs = attrs + " required"
		}
		if opt == "email" {
			inputType = TypeEmail
			continue
		}
		if opt == "uitextarea" {
			inputType = TypeTextarea
		}
		if opt == "uipassword" {
			inputType = TypePassword
		}
		for _, valOpt := range []string{"lenmin", "lenmax", "valmin", "valmax", "regexp"} {
			if strings.HasPrefix(opt, valOpt+":") {
				val := strings.Replace(opt, valOpt+":", "", 1)
				if valOpt == "regexp" {
					attrs = attrs + fmt.Sprintf(` pattern="%s"`, html.EscapeString(val))
					continue
				}

				i, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				switch valOpt {
				case "lenmin":
					attrs = attrs + fmt.Sprintf(` minlength="%d"`, i)
				case "lenmax":
					attrs = attrs + fmt.Sprintf(` maxlength="%d"`, i)
				case "valmin":
					attrs = attrs + fmt.Sprintf(` min="%d"`, i)
				case "valmax":
					attrs = attrs + fmt.Sprintf(` max="%d"`, i)
				}
			}
		}
	}

	return attrs, inputType
}
