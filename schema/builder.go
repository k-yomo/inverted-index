package schema

type Builder struct {
	Fields   []*FieldEntry `json:"fields"`
	FieldMap map[string]*FieldEntry
}

func NewBuilder() *Builder {
	return &Builder{
		FieldMap: make(map[string]*FieldEntry),
	}
}

func (b *Builder) AddTextField(fieldName string) {
	fieldEntry := newFieldEntry(len(b.Fields), fieldName, FieldTypeText)
	b.Fields = append(b.Fields, fieldEntry)
	b.FieldMap[fieldName] = fieldEntry
}

func (b *Builder) Build() *Schema {
	return &Schema{
		Fields:   b.Fields,
		fieldMap: b.FieldMap,
	}
}
