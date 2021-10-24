package schema

type FieldEntry struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	FieldType FieldType `json:"fieldType"`
}

func newFieldEntry(id int, name string, fieldType FieldType) *FieldEntry {
	return &FieldEntry{
		ID:        id,
		Name:      name,
		FieldType: fieldType,
	}
}
