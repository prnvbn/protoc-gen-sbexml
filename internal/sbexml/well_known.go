package sbexml

import "google.golang.org/protobuf/types/descriptorpb"

const (
	protoTimestampTypeName = ".google.protobuf.Timestamp"
	timestampTypeName      = "UTCTimestamp"
	timestampEpoch         = "unix"
	timestampTimeUnit      = "nanosecond"
)

var timestampPrimitiveType = primitiveType{
	Name:          timestampTypeName,
	PrimitiveType: "uint64",
}

func isTimestampField(field *descriptorpb.FieldDescriptorProto) bool {
	return field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE &&
		field.GetTypeName() == protoTimestampTypeName
}

func (g *generator) addTimestampType() {
	g.addPrimitive(timestampPrimitiveType)
}
