package form

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const FieldTypeNumber FieldType = "number"

type FieldNumber struct {
	FieldBase `yaml:",inline"`
	Min       *float64 `yaml:"min"`
	Max       *float64 `yaml:"max"`
}

func NewFieldNumber(name string) *FieldNumber {
	return &FieldNumber{
		FieldBase: FieldBase{
			Name: name,
			Type: FieldTypeNumber,
		},
	}
}

func init() {
	RegisterFieldType(FieldTypeNumber, func(n yaml.Node) (Field, error) {
		field := &FieldNumber{
			FieldBase: FieldBase{},
		}

		if err := n.Decode(field); err != nil {
			return nil, errors.WithStack(err)
		}

		return field, nil
	})
}
