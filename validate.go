package pom

import (
	"context"
	"regexp"
	"strconv"

	"github.com/go-openapi/spec"
)

func StringValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {
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
func IntegerValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {

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
func NumberValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {

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
func BooleanValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {
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
func MultiValidate(ctx context.Context, p spec.Parameter, v string, validate ...Validate) (Result, bool) {
	var rvi interface{}
	for _, vi := range validate {
		r, ok := vi(ctx, p, v)
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

var MobileRegexp = regexp.MustCompile("^((0\\d{2,3}-\\d{7,8})|(1[3|4|5|6|7|8|9][0-9]\\d{8}))$")

func MobileValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {
	i := MobileRegexp.FindString(v)
	ok := i != ""
	if !ok {
		return Result{
			ErrType: ErrType_Mobile,
		}, false
	}
	return Result{
		Value: v,
	}, true
}

var MailRegexp = regexp.MustCompile("^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$")

func MailValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {
	i := MailRegexp.FindString(v)
	ok := i != ""
	if !ok {
		return Result{
			ErrType: ErrType_Mail,
		}, false
	}
	return Result{
		Value: v,
	}, true
}

var IDCardRegexp = regexp.MustCompile("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)")

func IDCardValidate(ctx context.Context, p spec.Parameter, v string) (Result, bool) {
	i := IDCardRegexp.FindString(v)
	ok := i != ""
	if !ok {
		return Result{
			ErrType: ErrType_IDCard,
		}, false
	}
	return Result{
		Value: v,
	}, true
}
