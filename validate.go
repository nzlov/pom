package pom

import (
	"regexp"
	"strconv"

	"github.com/go-openapi/spec"
)

func StringValidate(p spec.Parameter, v string) (Result, bool) {
	r := []rune(v)
	l := int64(len(r))
	if p.MinLength != nil {
		if l < *p.MinLength {
			return Result{
				ErrType: ErrType_Range,
			}, false
		}
	}

	if p.MaxLength != nil {
		if l > *p.MaxLength {
			return Result{
				ErrType: ErrType_Range,
			}, false
		}
	}
	return Result{
		Value: v,
	}, true
}
func IntegerValidate(p spec.Parameter, v string) (Result, bool) {

	f, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return Result{
			ErrType: ErrType_Type,
		}, false
	}
	if p.Minimum != nil {
		if float64(f) < *p.Minimum {
			return Result{
				ErrType: ErrType_Range,
			}, false
		}
	}
	if p.Maximum != nil {
		if float64(f) > *p.Maximum {
			return Result{
				ErrType: ErrType_Range,
			}, false
		}
	}

	return Result{
		Value: f,
	}, true
}
func NumberValidate(p spec.Parameter, v string) (Result, bool) {

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return Result{
			ErrType: ErrType_Type,
		}, false
	}
	if p.Minimum != nil {
		if f < *p.Minimum {
			return Result{
				ErrType: ErrType_Range,
			}, false
		}
	}
	if p.Maximum != nil {
		if f > *p.Maximum {
			return Result{
				ErrType: ErrType_Range,
			}, false
		}
	}

	return Result{
		Value: f,
	}, true
}
func BooleanValidate(p spec.Parameter, v string) (Result, bool) {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return Result{
			ErrType: ErrType_Type,
		}, false
	}
	return Result{
		Value: b,
	}, true

}
func MultiValidate(p spec.Parameter, v string, validate ...Validate) (Result, bool) {
	var rvi interface{}
	for _, vi := range validate {
		r, ok := vi(p, v)
		if !ok {
			return r, false
		}
		if r.Value != nil {
			rvi = r.Value
		}
	}
	return Result{
		Value: rvi,
	}, true
}

var mobileRegexp = regexp.MustCompile("^((0\\d{2,3}-\\d{7,8})|(1[3|4|5|6|7|8|9][0-9]\\d{8}))$")

func MobileValidate(p spec.Parameter, v string) (Result, bool) {
	i := mobileRegexp.FindString(v)
	ok := i != ""
	if !ok {
		return Result{
			ErrType: ErrType_Format,
		}, false
	}
	return Result{
		Value: v,
	}, true
}

var mailRegexp = regexp.MustCompile("^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$")

func MailValidate(p spec.Parameter, v string) (Result, bool) {
	i := mailRegexp.FindString(v)
	ok := i != ""
	if !ok {
		return Result{
			ErrType: ErrType_Format,
		}, false
	}
	return Result{
		Value: v,
	}, true
}

var idCardRegexp = regexp.MustCompile("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)")

func IDCardValidate(p spec.Parameter, v string) (Result, bool) {
	i := idCardRegexp.FindString(v)
	ok := i != ""
	if !ok {
		return Result{
			ErrType: ErrType_Format,
		}, false
	}
	return Result{
		Value: v,
	}, true
}
