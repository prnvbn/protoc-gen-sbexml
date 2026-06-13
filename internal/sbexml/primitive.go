package sbexml

import "google.golang.org/protobuf/types/descriptorpb"

func (g *generator) addPrimitive(primitive primitiveType) {
	if _, ok := g.primitives[primitive.Name]; ok {
		return
	}
	g.primitives[primitive.Name] = struct{}{}
	g.schema.Types.Primitives = append(g.schema.Types.Primitives, primitive)
}

var protoPrimitiveTypes = map[descriptorpb.FieldDescriptorProto_Type]primitiveType{
	descriptorpb.FieldDescriptorProto_TYPE_DOUBLE: {
		Name:          "double",
		PrimitiveType: "double",
	},
	descriptorpb.FieldDescriptorProto_TYPE_FLOAT: {
		Name:          "float",
		PrimitiveType: "float",
	},
	descriptorpb.FieldDescriptorProto_TYPE_INT32: {
		Name:          "int32",
		PrimitiveType: "int32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_INT64: {
		Name:          "int64",
		PrimitiveType: "int64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_UINT32: {
		Name:          "uint32",
		PrimitiveType: "uint32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_UINT64: {
		Name:          "uint64",
		PrimitiveType: "uint64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_FIXED32: {
		Name:          "fixed32",
		PrimitiveType: "uint32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_FIXED64: {
		Name:          "fixed64",
		PrimitiveType: "uint64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_BOOL: {
		Name:          "bool",
		PrimitiveType: "uint8",
	},
	descriptorpb.FieldDescriptorProto_TYPE_STRING: {
		Name:              "string",
		PrimitiveType:     "char",
		CharacterEncoding: "UTF-8",
	},
	descriptorpb.FieldDescriptorProto_TYPE_BYTES: {
		Name:          "bytes",
		PrimitiveType: "uint8",
	},
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: {
		Name:          "sfixed32",
		PrimitiveType: "int32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: {
		Name:          "sfixed64",
		PrimitiveType: "int64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_SINT32: {
		Name:          "sint32",
		PrimitiveType: "int32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_SINT64: {
		Name:          "sint64",
		PrimitiveType: "int64",
	},
}
