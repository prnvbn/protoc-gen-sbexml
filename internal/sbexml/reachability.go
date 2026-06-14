package sbexml

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

type reachableTypes struct {
	enums           map[string]struct{}
	messages        map[string]struct{}
	pendingMessages []indexedMessage
}

func findReachableTypes(ti typeIndex, generatedFiles map[string]struct{}) (reachableTypes, error) {
	reachable := reachableTypes{
		enums:    map[string]struct{}{},
		messages: map[string]struct{}{},
	}

	for _, enum := range ti.enums {
		if _, ok := generatedFiles[enum.fileName]; ok {
			reachable.markEnum(enum.fullName)
		}
	}
	for _, message := range ti.messages {
		if _, ok := generatedFiles[message.fileName]; ok {
			reachable.markMessage(message)
		}
	}

	for len(reachable.pendingMessages) > 0 {
		message := reachable.pendingMessages[0]
		reachable.pendingMessages = reachable.pendingMessages[1:]

		for _, field := range message.descriptor.Field {
			if err := reachable.markFieldTypes(message, field, ti); err != nil {
				return reachableTypes{}, fmt.Errorf("resolve field %q on %s: %w", field.GetName(), message.name, err)
			}
		}
	}

	return reachable, nil
}

func (reachable *reachableTypes) markEnum(fullName string) {
	reachable.enums[fullName] = struct{}{}
}

func (reachable *reachableTypes) markMessage(im indexedMessage) {
	if _, ok := reachable.messages[im.fullName]; ok {
		return
	}
	reachable.messages[im.fullName] = struct{}{}
	reachable.pendingMessages = append(reachable.pendingMessages, im)
}

func (reachable *reachableTypes) markFieldTypes(im indexedMessage, field *descriptorpb.FieldDescriptorProto, ti typeIndex) error {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum, ok := ti.enumByName[field.GetTypeName()]
		if !ok {
			return fmt.Errorf("%w: %s for %s.%s", errUnsupportedEnumType, field.GetTypeName(), im.name, field.GetName())
		}
		reachable.markEnum(enum.fullName)
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if mapEntry, ok := ti.mapEntryMessages[field.GetTypeName()]; ok {
			for _, entryField := range mapEntry.descriptor.Field {
				if err := reachable.markFieldTypes(mapEntry, entryField, ti); err != nil {
					return fmt.Errorf("resolve map entry field %q: %w", entryField.GetName(), err)
				}
			}
			return nil
		}

		referenced, ok := ti.messageByName[field.GetTypeName()]
		if !ok {
			return fmt.Errorf("%w: %s for %s.%s", errUnsupportedMessageType, field.GetTypeName(), im.name, field.GetName())
		}
		reachable.markMessage(referenced)
	}

	return nil
}

func orderedEnums(ti typeIndex, fileToGenerate []string, generatedFiles map[string]struct{}, reachable map[string]struct{}) []indexedEnum {
	var result []indexedEnum
	seen := map[string]struct{}{}

	for _, fileName := range fileToGenerate {
		for _, enum := range ti.enums {
			if enum.fileName == fileName {
				result = appendReachableEnum(result, seen, reachable, enum)
			}
		}
	}
	for _, enum := range ti.enums {
		if _, ok := generatedFiles[enum.fileName]; ok {
			continue
		}
		result = appendReachableEnum(result, seen, reachable, enum)
	}

	return result
}

func appendReachableEnum(result []indexedEnum, seen map[string]struct{}, reachable map[string]struct{}, enum indexedEnum) []indexedEnum {
	if _, ok := reachable[enum.fullName]; !ok {
		return result
	}
	if _, ok := seen[enum.fullName]; ok {
		return result
	}
	seen[enum.fullName] = struct{}{}
	return append(result, enum)
}

func orderedMessages(ti typeIndex, fileToGenerate []string, generatedFiles map[string]struct{}, reachable map[string]struct{}) []indexedMessage {
	var result []indexedMessage
	seen := map[string]struct{}{}

	for _, fileName := range fileToGenerate {
		for _, message := range ti.messages {
			if message.fileName == fileName {
				result = appendReachableMessage(result, seen, reachable, message)
			}
		}
	}
	for _, message := range ti.messages {
		if _, ok := generatedFiles[message.fileName]; ok {
			continue
		}
		result = appendReachableMessage(result, seen, reachable, message)
	}

	return result
}

func appendReachableMessage(result []indexedMessage, seen map[string]struct{}, reachable map[string]struct{}, im indexedMessage) []indexedMessage {
	if _, ok := reachable[im.fullName]; !ok {
		return result
	}
	if _, ok := seen[im.fullName]; ok {
		return result
	}
	seen[im.fullName] = struct{}{}
	return append(result, im)
}
