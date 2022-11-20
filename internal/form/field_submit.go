package form

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const FieldTypeSubmit FieldType = "submit"

type FieldSubmit struct {
	FieldBase `yaml:",inline"`
}

func NewFieldSubmit(name string) *FieldEmail {
	return &FieldEmail{
		FieldBase: FieldBase{
			Name: name,
			Type: FieldTypeEmail,
		},
	}
}

func init() {
	RegisterFieldType(FieldTypeSubmit, func(n yaml.Node) (Field, error) {
		field := &FieldSubmit{
			FieldBase: FieldBase{},
		}

		if err := n.Decode(field); err != nil {
			return nil, errors.WithStack(err)
		}

		return field, nil
	})
}
