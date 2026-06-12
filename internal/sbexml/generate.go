package sbexml

import (
	"encoding/xml"
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const OutputFilename = "output.xml"

// Generate builds an SBE XML document for the files requested by protoc.
func Generate(request *pluginpb.CodeGeneratorRequest) (string, error) {
	g := newGenerator(request)
	if err := g.generate(); err != nil {
		return "", err
	}
	return g.xml()
}

type generator struct {
	request       *pluginpb.CodeGeneratorRequest
	schema        messageSchema
	files         map[string]*descriptorpb.FileDescriptorProto
	primitives    map[string]struct{}
	nextMessageID int
}

func newGenerator(request *pluginpb.CodeGeneratorRequest) *generator {
	files := make(map[string]*descriptorpb.FileDescriptorProto, len(request.ProtoFile))
	for _, file := range request.ProtoFile {
		files[file.GetName()] = file
	}

	return &generator{
		request:       request,
		schema:        newMessageSchema(),
		files:         files,
		primitives:    map[string]struct{}{},
		nextMessageID: 1,
	}
}

func (g *generator) generate() error {
	for _, name := range g.request.FileToGenerate {
		file, ok := g.files[name]
		if !ok {
			return fmt.Errorf("file to generate %q not found in request", name)
		}
		if err := g.addFile(file); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) addFile(file *descriptorpb.FileDescriptorProto) error {
	if g.schema.Package == "" {
		g.schema.Package = file.GetPackage()
	}

	enumNames := enumNames(file.GetPackage(), file.EnumType)
	for _, enum := range file.EnumType {
		g.schema.Types.Enums = append(g.schema.Types.Enums, buildEnum(enum))
	}

	for _, message := range file.MessageType {
		built, err := g.buildMessage(message, enumNames)
		if err != nil {
			return err
		}
		g.schema.Messages = append(g.schema.Messages, built)
		g.nextMessageID++
	}

	return nil
}

func (g *generator) buildMessage(message *descriptorpb.DescriptorProto, enumNames map[string]string) (messageType, error) {
	result := messageType{
		Name: message.GetName(),
		ID:   g.nextMessageID,
	}

	for _, field := range message.Field {
		fieldType, err := g.fieldType(message, field, enumNames)
		if err != nil {
			return messageType{}, err
		}
		result.Fields = append(result.Fields, fieldTypeXML{
			Name: field.GetName(),
			ID:   int(field.GetNumber()),
			Type: fieldType,
		})
	}

	return result, nil
}

func (g *generator) xml() (string, error) {
	data, err := xml.MarshalIndent(g.schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal XML: %w", err)
	}
	return xml.Header + string(data) + "\n", nil
}

func buildEnum(enum *descriptorpb.EnumDescriptorProto) enumType {
	result := enumType{
		Name:         enum.GetName(),
		EncodingType: "uint8",
	}

	for _, value := range enum.Value {
		result.Values = append(result.Values, validValue{
			Name:  value.GetName(),
			Value: value.GetNumber(),
		})
	}

	return result
}

func (g *generator) fieldType(message *descriptorpb.DescriptorProto, field *descriptorpb.FieldDescriptorProto, enumNames map[string]string) (string, error) {
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
		name, ok := enumNames[field.GetTypeName()]
		if !ok {
			return "", fmt.Errorf("unsupported enum type %s for %s.%s", field.GetTypeName(), message.GetName(), field.GetName())
		}
		return name, nil
	}

	primitive, ok := protoPrimitiveTypes[field.GetType()]
	if !ok {
		return "", fmt.Errorf("unsupported field type %s for %s.%s", field.GetType(), message.GetName(), field.GetName())
	}
	g.addPrimitive(primitive)
	return primitive.Name, nil
}

func (g *generator) addPrimitive(primitive primitiveType) {
	if _, ok := g.primitives[primitive.Name]; ok {
		return
	}
	g.primitives[primitive.Name] = struct{}{}
	g.schema.Types.Primitives = append(g.schema.Types.Primitives, primitive)
}

func enumNames(packageName string, enums []*descriptorpb.EnumDescriptorProto) map[string]string {
	names := make(map[string]string, len(enums))
	for _, enum := range enums {
		names[fullName(packageName, enum.GetName())] = enum.GetName()
	}
	return names
}

func fullName(packageName, name string) string {
	if packageName == "" {
		return "." + name
	}
	return "." + packageName + "." + name
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
