package structvalidator

import (
	"reflect"
	"regexp"
	"strings"
)

type ValueValidation struct {
	LenMin int
	LenMax int
	ValMin int64
	ValMax int64
	Regexp *regexp.Regexp
	Flags  int64
}

// values used with flags
const (
	_            = iota
	ValMinNotNil = 1 << iota
	ValMaxNotNil
	Required
	Email
)

func (v *ValueValidation) ValidateReflectValue(value reflect.Value) (ok bool, failureFlags int) {
	minCanBeZero := false
	maxCanBeZero := false
	if v.Flags&ValMinNotNil > 0 {
		minCanBeZero = true
	}
	if v.Flags&ValMaxNotNil > 0 {
		maxCanBeZero = true
	}

	if v.Flags&Required > 0 {
		if value.Type().Name() == "string" && value.String() == "" {
			return false, FailEmpty
		}
		if strings.HasPrefix(value.Type().Name(), "int") && value.Int() == 0 && !minCanBeZero && !maxCanBeZero && v.ValMin == 0 && v.ValMax == 0 {
			return false, FailZero
		}
	}

	if value.Type().Name() == "string" {
		if v.LenMin > 0 && len(value.String()) < v.LenMin {
			return false, FailLenMin
		}
		if v.LenMax > 0 && len(value.String()) > v.LenMax {
			return false, FailLenMax
		}

		if v.Regexp != nil {
			if !v.Regexp.MatchString(value.String()) {
				return false, FailRegexp
			}
		}

		if v.Flags&Email > 0 {
			var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if !emailRegex.MatchString(value.String()) {
				return false, FailEmail
			}
		}
	}

	if strings.HasPrefix(value.Type().Name(), "int") {
		if (v.ValMin != 0 || minCanBeZero) && v.ValMin > value.Int() {
			return false, FailValMin
		}
		if (v.ValMax != 0 || maxCanBeZero) && v.ValMax < value.Int() {
			return false, FailValMax
		}
	}

	return true, 0
}

func NewValueValidation() *ValueValidation {
	return &ValueValidation{
		LenMin: -1,
		LenMax: -1,
	}
}
