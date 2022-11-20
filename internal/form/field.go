package form

type FieldType string

const (
	FieldTypeDateTimeLocal FieldType = "datetime-local"
	FieldTypeMonth         FieldType = "month"
	FieldTypeTime          FieldType = "time"
	FieldTypeWeek          FieldType = "week"
	FieldTypeColor         FieldType = "color"
	FieldTypeRange         FieldType = "range"
	FieldTypePhone         FieldType = "phone"
	FieldTypeURL           FieldType = "url"
)

type Field interface {
	FieldName() string
	FieldType() FieldType
}
