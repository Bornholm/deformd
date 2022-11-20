package form

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const FieldTypeEmail FieldType = "email"

type FieldEmail struct {
	FieldBase `yaml:",inline"`
}

func NewFieldEmail(name string) *FieldEmail {
	return &FieldEmail{
		FieldBase: FieldBase{
			Name: name,
			Type: FieldTypeEmail,
		},
	}
}

func init() {
	RegisterFieldType(FieldTypeEmail, func(n yaml.Node) (Field, error) {
		field := &FieldEmail{
			FieldBase: FieldBase{},
		}

		if err := n.Decode(field); err != nil {
			return nil, errors.WithStack(err)
		}

		return field, nil
	})
}
