package schema

type Field uint32

type FieldAndFieldValues struct {
	Field       Field
	FieldValues []*FieldValue
}

type FieldValue struct {
	Field Field
	Value Value
}
