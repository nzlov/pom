package pom // import "github.com/nzlov/pom"

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

var def *Pom

func Default() *Pom {
	if def == nil {
		def = New()
	}
	return def
}

// 验证成功返回时Result.Value赋值正确的数据类型
type Validate func(context.Context, spec.Parameter, string) (Result, bool)

func RegisterValidate(name string, f Validate) {
	Default().RegisterValidate(name, f)
}
func RegisterModelValidate(name string, f Validate) {
	Default().RegisterModelValidate(name, f)
}
func Action(ctx context.Context, id string, props Props) (Result, bool) {
	return Default().Action(ctx, id, props)
}
func Parse(data []byte) error {
	return Default().Parse(data)
}

type Pom struct {
	actions             map[string]*spec.Operation
	security            spec.SecurityDefinitions
	customValidate      map[string]Validate
	customModelValidate map[string]Validate

	kF func(string) string
}

func New() *Pom {
	return &Pom{
		actions:             map[string]*spec.Operation{},
		customValidate:      map[string]Validate{},
		customModelValidate: map[string]Validate{},
	}
}

func NewWithData(data []byte) (*Pom, error) {
	p := New()
	err := p.Parse(data)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (p *Pom) SetKeyFunc(f func(string) string) {
	p.kF = f
}

func (p *Pom) RegisterValidate(name string, f Validate) {
	p.customValidate[name] = f
}
func (p *Pom) RegisterModelValidate(name string, f Validate) {
	logrus.Debugln("Pom RegisterModelValidate:", name)
	p.customModelValidate[name] = f
}
func (p *Pom) Parse(data []byte) error {
	swagger := spec.Swagger{}
	err := swagger.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	p.security = swagger.SecurityDefinitions

	for k, v := range swagger.Paths.Paths {
		k = p.kF(k)
		if v.Get != nil {
			logrus.Debugln("Pom Parse:Get:", k)
			p.actions["GET "+k] = v.Get
		}
		if v.Post != nil {
			logrus.Debugln("Pom Parse:Post:", k)
			p.actions["POST "+k] = v.Post
		}
		if v.Delete != nil {
			logrus.Debugln("Pom Parse:Delete:", k)
			p.actions["DELETE "+k] = v.Delete
		}
	}
	return nil
}

type Props interface {
	Get(spec.Parameter) string
	Set(string, interface{})
	Security(*spec.Operation) bool
}

type ErrType string

const (
	ErrType_Action     ErrType = "no action"
	ErrType_Permission         = "no permission"
	ErrType_Format             = "format error"
	ErrType_Type               = "type error"
	ErrType_Required           = "required"
	ErrType_Range              = "range error"

	ErrType_Mobile = "mobile"
	ErrType_IDCard = "idcard"
	ErrType_Mail   = "mail"
)

type Result struct {
	ID        string
	Name      string
	ErrType   ErrType
	Value     interface{} `json:"-"`
	MinLength *int64
	MaxLength *int64
	Minimum   *float64
	Maximum   *float64
	Enum      []interface{}
}

// Action 验证id对应的Props是否通过，如果通过则通过Props.Set设置属性到Props上，并返回true
func (p *Pom) Action(ctx context.Context, id string, props Props) (Result, bool) {
	now := time.Now()
	defer func() {
		logrus.Debugln("Pom Action:", id, time.Since(now))
	}()
	action, ok := p.actions[id]
	if !ok {
		return Result{
			Name:    id,
			ErrType: ErrType_Action,
		}, false
	}
	if !props.Security(action) {
		return Result{
			Name:    id,
			ErrType: ErrType_Permission,
		}, false
	}
	for _, v := range action.Parameters {
		//logrus.Debugln("Pom Action:", id, "Param:", v.Name)
		var validate Validate
		if f, ok := p.customValidate[v.Name]; ok {
			validate = f
		} else {
			switch v.Type {
			case "string":
				validate = StringValidate
			case "integer":
				validate = IntegerValidate
			case "number":
				validate = NumberValidate
			case "boolean":
				validate = BooleanValidate
			default:
				if f, ok := p.customModelValidate[v.Type]; ok {
					validate = f
				} else {
					logrus.Errorln("Pom Action No Custom Model Validate:", v.Type)
					continue
				}
			}
		}

		pv := props.Get(v)
		if pv == "" {
			if v.Required {
				return Result{
					Name:    v.Name,
					ErrType: ErrType_Required,
				}, false
			} else {
				continue
			}
		}

		if len(v.Enum) > 0 {
			has := false
			for _, e := range v.Enum {
				if pv == fmt.Sprint(e) {
					has = true
					break
				}

			}
			if !has {
				return Result{
					Name:      v.Name,
					ErrType:   ErrType_Range,
					MinLength: v.MinLength,
					MaxLength: v.MaxLength,
					Minimum:   v.Minimum,
					Maximum:   v.Maximum,
					Enum:      v.Enum,
				}, false
			}
		}

		r, ok := validate(ctx, v, pv)
		if !ok {
			r.Name = v.Name
			r.MinLength = v.MinLength
			r.MaxLength = v.MaxLength
			r.Minimum = v.Minimum
			r.Maximum = v.Maximum
			r.Enum = v.Enum
			return r, false
		}
		props.Set(v.Name, r.Value)
	}

	return Result{ID: action.ID}, true
}
