package sbexml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const OutputFilename = "output.xml"

var (
	errFileToGenerateNotFound = errors.New("file to generate not found")
	errUnsupportedEnumType    = errors.New("unsupported enum type")
	errUnsupportedMessageType = errors.New("unsupported message type")
	errUnsupportedFieldType   = errors.New("unsupported field type")
)

// Generate builds an SBE XML document for the files requested by protoc.
func Generate(request *pluginpb.CodeGeneratorRequest) (string, error) {
	g := newGenerator(request)
	if err := g.generate(); err != nil {
		return "", fmt.Errorf("generate schema: %w", err)
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
			return fmt.Errorf("%w: %q", errFileToGenerateNotFound, name)
		}
		if err := g.addFile(file); err != nil {
			return fmt.Errorf("add file %q: %w", name, err)
		}
	}
	return nil
}

func (g *generator) addFile(file *descriptorpb.FileDescriptorProto) error {
	if g.schema.Package == "" {
		g.schema.Package = file.GetPackage()
	}

	types := indexFileTypes(file)
	for _, enum := range types.enums {
		g.schema.Types.Enums = append(g.schema.Types.Enums, buildEnum(enum))
	}

	for _, message := range types.messages {
		built, err := g.buildMessage(message, types)
		if err != nil {
			return fmt.Errorf("build message %q: %w", message.name, err)
		}
		g.schema.Messages = append(g.schema.Messages, built)
		g.nextMessageID++
	}

	return nil
}

func (g *generator) buildMessage(message indexedMessage, types fileTypes) (messageType, error) {
	result := messageType{
		Name: message.name,
		ID:   g.nextMessageID,
	}

	for _, field := range message.descriptor.Field {
		fieldType, err := g.fieldType(message, field, types)
		if err != nil {
			return messageType{}, fmt.Errorf("resolve field %q: %w", field.GetName(), err)
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

func buildEnum(enum indexedEnum) enumType {
	result := enumType{
		Name:         enum.name,
		EncodingType: "uint8",
	}

	for _, value := range enum.descriptor.Value {
		result.Values = append(result.Values, validValue{
			Name:  value.GetName(),
			Value: value.GetNumber(),
		})
	}

	return result
}

func (g *generator) fieldType(message indexedMessage, field *descriptorpb.FieldDescriptorProto, types fileTypes) (string, error) {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		name, ok := types.enumNames[field.GetTypeName()]
		if !ok {
			return "", fmt.Errorf("%w: %s for %s.%s", errUnsupportedEnumType, field.GetTypeName(), message.name, field.GetName())
		}
		return name, nil
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		name, ok := types.messageNames[field.GetTypeName()]
		if !ok {
			return "", fmt.Errorf("%w: %s for %s.%s", errUnsupportedMessageType, field.GetTypeName(), message.name, field.GetName())
		}
		return name, nil
	default:
		primitive, ok := protoPrimitiveTypes[field.GetType()]
		if !ok {
			return "", fmt.Errorf("%w: %s for %s.%s", errUnsupportedFieldType, field.GetType(), message.name, field.GetName())
		}
		g.addPrimitive(primitive)
		return primitive.Name, nil
	}
}

func (g *generator) addPrimitive(primitive primitiveType) {
	if _, ok := g.primitives[primitive.Name]; ok {
		return
	}
	g.primitives[primitive.Name] = struct{}{}
	g.schema.Types.Primitives = append(g.schema.Types.Primitives, primitive)
}

type fileTypes struct {
	enumNames    map[string]string
	messageNames map[string]string
	enums        []indexedEnum
	messages     []indexedMessage
}

type indexedEnum struct {
	descriptor *descriptorpb.EnumDescriptorProto
	name       string
}

type indexedMessage struct {
	descriptor *descriptorpb.DescriptorProto
	name       string
}

func indexFileTypes(file *descriptorpb.FileDescriptorProto) fileTypes {
	types := fileTypes{
		enumNames:    map[string]string{},
		messageNames: map[string]string{},
	}

	for _, enum := range file.EnumType {
		types.addEnum(file.GetPackage(), nil, enum)
	}
	for _, message := range file.MessageType {
		types.addMessage(file.GetPackage(), nil, message)
	}

	return types
}

func (types *fileTypes) addEnum(packageName string, parentPath []string, enum *descriptorpb.EnumDescriptorProto) {
	path := appendPath(parentPath, enum.GetName())
	name := xmlTypeName(path)

	types.enumNames[protoFullName(packageName, path)] = name
	types.enums = append(types.enums, indexedEnum{
		descriptor: enum,
		name:       name,
	})
}

func (types *fileTypes) addMessage(packageName string, parentPath []string, message *descriptorpb.DescriptorProto) {
	path := appendPath(parentPath, message.GetName())
	name := xmlTypeName(path)

	types.messageNames[protoFullName(packageName, path)] = name
	types.messages = append(types.messages, indexedMessage{
		descriptor: message,
		name:       name,
	})

	for _, enum := range message.EnumType {
		types.addEnum(packageName, path, enum)
	}
	for _, nested := range message.NestedType {
		types.addMessage(packageName, path, nested)
	}
}

func appendPath(path []string, name string) []string {
	next := make([]string, 0, len(path)+1)
	next = append(next, path...)
	return append(next, name)
}

func xmlTypeName(path []string) string {
	return strings.Join(path, "_")
}

func protoFullName(packageName string, path []string) string {
	name := strings.Join(path, ".")
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
