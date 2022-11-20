package config

import (
	"github.com/Bornholm/deformd/internal/form"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type FormFields []form.Field

type FormConfig struct {
	Title   InterpolatedString `yaml:"title"`
	Fields  FormFields         `yaml:"fields"`
	Handler Handler            `yaml:"handler"`
}

type typedField struct {
	Type form.FieldType `yaml:type`
}

func (f *FormFields) UnmarshalYAML(value *yaml.Node) error {
	nodes := make([]yaml.Node, 0)

	if err := value.Decode(&nodes); err != nil {
		return errors.WithStack(err)
	}

	for _, n := range nodes {
		tf := typedField{}

		if err := n.Decode(&tf); err != nil {
			return errors.WithStack(err)
		}

		field, err := form.CreateField(tf.Type, n)
		if err != nil {
			return errors.WithStack(err)
		}

		*f = append(*f, field)
	}

	return nil
}

func NewDefaultFormConfig() FormConfig {
	return FormConfig{
		Fields: make([]form.Field, 0),
	}
}
