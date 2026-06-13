package sbexml

import "google.golang.org/protobuf/types/descriptorpb"

const (
	groupSizeEncodingType = "groupSizeEncoding"
	repeatedValueFieldID  = 1
	repeatedValueName     = "value"
)

func (g *generator) buildGroup(message indexedMessage, field *descriptorpb.FieldDescriptorProto, types fileTypes) (groupType, error) {
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
