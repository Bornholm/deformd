package form

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var ErrFieldTypeNotRegistered = errors.New("field type not registered")

var defaultFieldTypeRegistry = NewFieldTypeRegistry()

func RegisterFieldType(fieldType FieldType, unmarshalFunc FieldUnmarshalFunc) {
	defaultFieldTypeRegistry.Register(fieldType, unmarshalFunc)
}

func CreateField(fieldType FieldType, node yaml.Node) (Field, error) {
	field, err := defaultFieldTypeRegistry.Create(fieldType, node)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return field, nil
}

type FieldUnmarshalFunc func(n yaml.Node) (Field, error)

type FieldTypeRegistry struct {
	fieldTypes map[FieldType]FieldUnmarshalFunc
}

func (r *FieldTypeRegistry) Register(fieldType FieldType, unmarshalFunc FieldUnmarshalFunc) {
	r.fieldTypes[fieldType] = unmarshalFunc
}

func (r *FieldTypeRegistry) Create(fieldType FieldType, n yaml.Node) (Field, error) {
	unmarshall, exists := r.fieldTypes[fieldType]
	if !exists {
		return nil, errors.Wrapf(ErrFieldTypeNotRegistered, "could not create field type '%s'", fieldType)
	}

	field, err := unmarshall(n)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return field, nil
}

func NewFieldTypeRegistry() *FieldTypeRegistry {
	return &FieldTypeRegistry{
		fieldTypes: make(map[FieldType]FieldUnmarshalFunc),
	}
}
