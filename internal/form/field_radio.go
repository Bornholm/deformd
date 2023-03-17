package form

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const FieldTypeRadio FieldType = "radio"

type FieldRadio struct {
	FieldBase `yaml:",inline"`
	Values    []FieldRadioValue `yaml:"values"`
}

type FieldRadioValue struct {
	Label string
	Value string
}

func NewFieldRadio(name string) *FieldRadio {
	return &FieldRadio{
		FieldBase: FieldBase{
			Name: name,
			Type: FieldTypeRadio,
		},
		Values: make([]FieldRadioValue, 0),
	}
}

func init() {
	RegisterFieldType(FieldTypeRadio, func(n yaml.Node) (Field, error) {
		field := &FieldRadio{
			FieldBase: FieldBase{},
		}

		if err := n.Decode(field); err != nil {
			return nil, errors.WithStack(err)
		}

		return field, nil
	})
}
