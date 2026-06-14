package sbexml

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

type oneofIndex struct {
	byIndex map[int]oneofFields
	seen    map[int]struct{}
}

type oneofFields struct {
	index  int
	name   string
	fields []*descriptorpb.FieldDescriptorProto
}

func indexOneofs(message *descriptorpb.DescriptorProto) oneofIndex {
	oneofs := oneofIndex{
		byIndex: map[int]oneofFields{},
		seen:    map[int]struct{}{},
	}

	for index, declaration := range message.OneofDecl {
		oneofs.byIndex[index] = oneofFields{
			index: index,
			name:  declaration.GetName(),
		}
	}

	for _, field := range message.Field {
		if !isOneofField(field) {
			continue
		}
		index := int(field.GetOneofIndex())
		oneof, ok := oneofs.byIndex[index]
		if !ok {
			continue
		}
		oneof.fields = append(oneof.fields, field)
		oneofs.byIndex[index] = oneof
	}

	return oneofs
}

func isOneofField(field *descriptorpb.FieldDescriptorProto) bool {
	return field.OneofIndex != nil && !field.GetProto3Optional()
}

func (oneofs oneofIndex) emitted(index int) bool {
	if _, ok := oneofs.seen[index]; ok {
		return true
	}
	oneofs.seen[index] = struct{}{}
	return false
}

func (g *generator) buildOneof(im indexedMessage, oneof oneofFields, ti typeIndex) (fieldTypeXML, error) {
	compositeName := xmlTypeName([]string{im.name, oneof.name})
	fields := make([]fieldTypeXML, 0, len(oneof.fields))

	for _, field := range oneof.fields {
		fieldType, err := g.fieldType(im, field, ti)
		if err != nil {
			return fieldTypeXML{}, fmt.Errorf("resolve option %q: %w", field.GetName(), err)
		}
		fields = append(fields, fieldTypeXML{
			Name: field.GetName(),
			ID:   int(field.GetNumber()),
			Type: fieldType,
		})
	}

	g.addOneofComposite(compositeName, fields)
	return fieldTypeXML{
		Name: oneof.name,
		ID:   int(oneof.fields[0].GetNumber()),
		Type: compositeName,
	}, nil
}

func (g *generator) addOneofComposite(name string, fields []fieldTypeXML) {
	if _, ok := g.composites[name]; ok {
		return
	}
	g.composites[name] = struct{}{}
	g.schema.Types.Composites = append(g.schema.Types.Composites, compositeType{
		Name:   name,
		Fields: fields,
	})
}
