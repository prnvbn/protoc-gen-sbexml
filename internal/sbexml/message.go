package sbexml

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) buildMessage(im indexedMessage, ti typeIndex) (messageType, error) {
	result := messageType{
		Name: im.name,
		ID:   g.nextMessageID,
	}
	oneofs := indexOneofs(im.descriptor)

	for _, field := range im.descriptor.Field {
		if isOneofField(field) {
			oneof, ok := oneofs.byIndex[int(field.GetOneofIndex())]
			if !ok {
				return messageType{}, fmt.Errorf("%w: %s.%s", errOneofDeclarationNotFound, im.name, field.GetName())
			}
			if oneofs.emitted(oneof.index) {
				continue
			}

			built, err := g.buildOneof(im, oneof, ti)
			if err != nil {
				return messageType{}, fmt.Errorf("resolve oneof %q: %w", oneof.name, err)
			}
			result.Fields = append(result.Fields, built)
			continue
		}

		if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			group, err := g.buildGroup(im, field, ti)
			if err != nil {
				return messageType{}, fmt.Errorf("resolve repeated field %q: %w", field.GetName(), err)
			}
			result.Groups = append(result.Groups, group)
			continue
		}

		fieldType, err := g.fieldType(im, field, ti)
		if err != nil {
			return messageType{}, fmt.Errorf("resolve field %q: %w", field.GetName(), err)
		}
		result.Fields = append(result.Fields, fieldTypeXML{
			Name:     field.GetName(),
			ID:       int(field.GetNumber()),
			Type:     fieldType,
			Presence: fieldPresence(field),
		})
	}

	return result, nil
}

func fieldPresence(field *descriptorpb.FieldDescriptorProto) string {
	if field.GetProto3Optional() {
		return "optional"
	}
	return ""
}

func (g *generator) fieldType(im indexedMessage, field *descriptorpb.FieldDescriptorProto, ti typeIndex) (string, error) {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum, ok := ti.enumByName[field.GetTypeName()]
		if !ok {
			return "", fmt.Errorf("%w: %s for %s.%s", errUnsupportedEnumType, field.GetTypeName(), im.name, field.GetName())
		}
		return enum.name, nil
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		referenced, ok := ti.messageByName[field.GetTypeName()]
		if !ok {
			return "", fmt.Errorf("%w: %s for %s.%s", errUnsupportedMessageType, field.GetTypeName(), im.name, field.GetName())
		}
		return referenced.name, nil
	default:
		primitive, ok := protoPrimitiveTypes[field.GetType()]
		if !ok {
			return "", fmt.Errorf("%w: %s for %s.%s", errUnsupportedFieldType, field.GetType(), im.name, field.GetName())
		}
		g.addPrimitive(primitive)
		return primitive.Name, nil
	}
}
