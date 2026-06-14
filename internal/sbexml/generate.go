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
	composites    map[string]struct{}
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
		composites:    map[string]struct{}{},
		nextMessageID: 1,
	}
}

func (g *generator) generate() error {
	ti := indexTypes(g.request.ProtoFile)
	generatedFiles, err := g.generatedFiles()
	if err != nil {
		return err
	}

	reachable, err := findReachableTypes(ti, generatedFiles)
	if err != nil {
		return fmt.Errorf("find reachable types: %w", err)
	}

	g.addEnums(orderedEnums(ti, g.request.FileToGenerate, generatedFiles, reachable.enums))

	if err := g.addMessages(orderedMessages(ti, g.request.FileToGenerate, generatedFiles, reachable.messages), ti); err != nil {
		return err
	}

	return nil
}

func (g *generator) generatedFiles() (map[string]struct{}, error) {
	generatedFiles := map[string]struct{}{}

	for _, name := range g.request.FileToGenerate {
		file, ok := g.files[name]
		if !ok {
			return nil, fmt.Errorf("%w: %q", errFileToGenerateNotFound, name)
		}
		if g.schema.Package == "" {
			g.schema.Package = file.GetPackage()
		}
		generatedFiles[name] = struct{}{}
	}

	return generatedFiles, nil
}

func (g *generator) addEnums(enums []indexedEnum) {
	for _, enum := range enums {
		g.schema.Types.Enums = append(g.schema.Types.Enums, buildEnum(enum))
	}
}

func (g *generator) addMessages(messages []indexedMessage, ti typeIndex) error {
	for _, message := range messages {
		built, err := g.buildMessage(message, ti)
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
