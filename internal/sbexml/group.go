package sbexml

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	groupSizeEncodingType = "groupSizeEncoding"
	repeatedValueFieldID  = 1
	repeatedValueName     = "value"
)

func (g *generator) buildGroup(message indexedMessage, field *descriptorpb.FieldDescriptorProto, types fileTypes) (groupType, error) {
	if mapEntry, ok := types.mapEntryMessages[field.GetTypeName()]; ok {
		return g.buildMapGroup(field, mapEntry, types)
	}

	fieldType, err := g.fieldType(message, field, types)
	if err != nil {
		return groupType{}, err
	}
	g.addGroupSizeEncoding()

	return groupType{
		Name:          field.GetName(),
		ID:            int(field.GetNumber()),
		DimensionType: groupSizeEncodingType,
		Fields: []fieldTypeXML{
			{
				Name: repeatedValueName,
				ID:   repeatedValueFieldID,
				Type: fieldType,
			},
		},
	}, nil
}

func (g *generator) buildMapGroup(field *descriptorpb.FieldDescriptorProto, mapEntry indexedMessage, types fileTypes) (groupType, error) {
	fields := make([]fieldTypeXML, 0, len(mapEntry.descriptor.Field))
	for _, entryField := range mapEntry.descriptor.Field {
		fieldType, err := g.fieldType(mapEntry, entryField, types)
		if err != nil {
			return groupType{}, fmt.Errorf("resolve map entry field %q: %w", entryField.GetName(), err)
		}
		fields = append(fields, fieldTypeXML{
			Name: entryField.GetName(),
			ID:   int(entryField.GetNumber()),
			Type: fieldType,
		})
	}

	g.addGroupSizeEncoding()

	return groupType{
		Name:          field.GetName(),
		ID:            int(field.GetNumber()),
		DimensionType: groupSizeEncodingType,
		Fields:        fields,
	}, nil
}

func (g *generator) addGroupSizeEncoding() {
	if _, ok := g.composites[groupSizeEncodingType]; ok {
		return
	}
	g.composites[groupSizeEncodingType] = struct{}{}
	g.schema.Types.Composites = append(g.schema.Types.Composites, compositeType{
		Name: groupSizeEncodingType,
		Types: []primitiveType{
			{
				Name:          "blockLength",
				PrimitiveType: "uint16",
			},
			{
				Name:          "numInGroup",
				PrimitiveType: "uint16",
			},
		},
	})
}
