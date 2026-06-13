package sbexml

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) buildMessage(message indexedMessage, types fileTypes) (messageType, error) {
	result := messageType{
		Name: message.name,
		ID:   g.nextMessageID,
	}
	oneofs := indexOneofs(message.descriptor)

	for _, field := range message.descriptor.Field {
		if isOneofField(field) {
			oneof, ok := oneofs.byIndex[int(field.GetOneofIndex())]
			if !ok {
				return messageType{}, fmt.Errorf("%w: %s.%s", errOneofDeclarationNotFound, message.name, field.GetName())
			}
			if oneofs.emitted(oneof.index) {
				continue
			}

			built, err := g.buildOneof(message, oneof, types)
			if err != nil {
				return messageType{}, fmt.Errorf("resolve oneof %q: %w", oneof.name, err)
			}
			result.Fields = append(result.Fields, built)
			continue
		}

		if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			group, err := g.buildGroup(message, field, types)
			if err != nil {
				return messageType{}, fmt.Errorf("resolve repeated field %q: %w", field.GetName(), err)
			}
			result.Groups = append(result.Groups, group)
			continue
		}

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
