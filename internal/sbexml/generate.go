package sbexml

import (
	"encoding/xml"
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const OutputFilename = "output.xml"

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
	groups        map[string]struct{}
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
		groups:        map[string]struct{}{},
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

func (g *generator) xml() (string, error) {
	data, err := xml.MarshalIndent(g.schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal XML: %w", err)
	}
	return xml.Header + string(data) + "\n", nil
}
