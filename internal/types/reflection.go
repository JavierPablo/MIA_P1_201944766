package types

type SystemType interface{
	total_size() int32;
	name() string
}
type StructType interface{
	SystemType
	fields() []Field
}

type Field struct{
	concrete_type SystemType
	label string
}

type StructInfo struct{
	total_size int32
}