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

func (g *generator) buildGroup(im indexedMessage, field *descriptorpb.FieldDescriptorProto, ti typeIndex) (groupType, error) {
	if mapEntry, ok := ti.mapEntryMessages[field.GetTypeName()]; ok {
		return g.buildMapGroup(field, mapEntry, ti)
	}

	valueField, err := g.buildField(im, field, repeatedValueName, repeatedValueFieldID, ti)
	if err != nil {
		return groupType{}, err
	}
	g.addGroupSizeEncoding()

	return groupType{
		Name:          field.GetName(),
		ID:            int(field.GetNumber()),
		DimensionType: groupSizeEncodingType,
		Fields:        []fieldTypeXML{valueField},
	}, nil
}

func (g *generator) buildMapGroup(field *descriptorpb.FieldDescriptorProto, mapEntry indexedMessage, ti typeIndex) (groupType, error) {
	fields := make([]fieldTypeXML, 0, len(mapEntry.descriptor.Field))
	for _, entryField := range mapEntry.descriptor.Field {
		built, err := g.buildField(mapEntry, entryField, entryField.GetName(), int(entryField.GetNumber()), ti)
		if err != nil {
			return groupType{}, fmt.Errorf("resolve map entry field %q: %w", entryField.GetName(), err)
		}
		fields = append(fields, built)
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
