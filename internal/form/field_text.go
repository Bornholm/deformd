package form

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const FieldTypeText FieldType = "text"

type FieldText struct {
	FieldBase `yaml:",inline"`
	Multiline *bool `yaml:"multiline"`
}

func NewFieldText(name string) *FieldText {
	return &FieldText{
		FieldBase: FieldBase{
			Name: name,
			Type: FieldTypeText,
		},
	}
}

func init() {
	RegisterFieldType(FieldTypeText, func(n yaml.Node) (Field, error) {
		field := &FieldText{
			FieldBase: FieldBase{},
		}

		if err := n.Decode(field); err != nil {
			return nil, errors.WithStack(err)
		}

		return field, nil
	})
}
