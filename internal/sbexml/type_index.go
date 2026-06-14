package sbexml

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type typeIndex struct {
	enums            []indexedEnum
	messages         []indexedMessage
	enumByName       map[string]indexedEnum
	messageByName    map[string]indexedMessage
	mapEntryMessages map[string]indexedMessage
}

type indexedEnum struct {
	descriptor *descriptorpb.EnumDescriptorProto
	name       string
	fullName   string
	fileName   string
}

type indexedMessage struct {
	descriptor *descriptorpb.DescriptorProto
	name       string
	fullName   string
	fileName   string
}

func indexTypes(files []*descriptorpb.FileDescriptorProto) typeIndex {
	ti := typeIndex{
		enumByName:       map[string]indexedEnum{},
		messageByName:    map[string]indexedMessage{},
		mapEntryMessages: map[string]indexedMessage{},
	}

	for _, file := range files {
		ti.addFile(file)
	}

	return ti
}

func (ti *typeIndex) addFile(file *descriptorpb.FileDescriptorProto) {
	for _, enum := range file.EnumType {
		ti.addEnum(file.GetName(), file.GetPackage(), nil, enum)
	}
	for _, message := range file.MessageType {
		ti.addMessage(file.GetName(), file.GetPackage(), nil, message)
	}
}

func (ti *typeIndex) addEnum(fileName string, packageName string, parentPath []string, enum *descriptorpb.EnumDescriptorProto) {
	path := appendPath(parentPath, enum.GetName())
	name := xmlTypeName(path)
	fullName := protoFullName(packageName, path)

	indexed := indexedEnum{
		descriptor: enum,
		name:       name,
		fullName:   fullName,
		fileName:   fileName,
	}
	ti.enumByName[fullName] = indexed
	ti.enums = append(ti.enums, indexed)
}

func (ti *typeIndex) addMessage(fileName string, packageName string, parentPath []string, message *descriptorpb.DescriptorProto) {
	path := appendPath(parentPath, message.GetName())
	name := xmlTypeName(path)
	fullName := protoFullName(packageName, path)

	indexed := indexedMessage{
		descriptor: message,
		name:       name,
		fullName:   fullName,
		fileName:   fileName,
	}

	if isMapEntry(message) {
		ti.mapEntryMessages[fullName] = indexed
		return
	}

	ti.messageByName[fullName] = indexed
	ti.messages = append(ti.messages, indexed)

	for _, enum := range message.EnumType {
		ti.addEnum(fileName, packageName, path, enum)
	}
	for _, nested := range message.NestedType {
		ti.addMessage(fileName, packageName, path, nested)
	}
}

func isMapEntry(message *descriptorpb.DescriptorProto) bool {
	return message.GetOptions().GetMapEntry()
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
