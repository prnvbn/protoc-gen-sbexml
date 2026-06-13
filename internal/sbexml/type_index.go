package sbexml

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

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
