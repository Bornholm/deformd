package form

type FieldBase struct {
	Label       string    `yaml:"label"`
	Name        string    `yaml:"name"`
	Type        FieldType `yaml:"type"`
	Placeholder string    `yaml:"placeholder"`
	Required    *bool     `yaml:"required"`
}

func (f *FieldBase) FieldName() string {
	return f.Name
}

func (f *FieldBase) FieldType() FieldType {
	return f.Type
}
