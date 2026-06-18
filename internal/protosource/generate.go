package protosource

import (
	"context"
	"fmt"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/prnvbn/protoc-gen-sbexml/internal/sbexml"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const inputFilename = "input.proto"

// GenerateSBE compiles one protobuf source file and generates SBE XML from it.
func GenerateSBE(ctx context.Context, source string) (string, error) {
	if strings.TrimSpace(source) == "" {
		return "", fmt.Errorf("proto input is empty")
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			Accessor: protocompile.SourceAccessorFromMap(map[string]string{
				inputFilename: source,
			}),
		}),
	}

	files, err := compiler.Compile(ctx, inputFilename)
	if err != nil {
		return "", fmt.Errorf("compile proto: %w", err)
	}
	if len(files) != 1 {
		return "", fmt.Errorf("compile proto: expected 1 file, got %d", len(files))
	}

	protoFiles := collectFileDescriptors(files[0])
	return sbexml.Generate(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{inputFilename},
		ProtoFile:      protoFiles,
	})
}

func collectFileDescriptors(file protoreflect.FileDescriptor) []*descriptorpb.FileDescriptorProto {
	seen := map[string]struct{}{}
	var descriptors []*descriptorpb.FileDescriptorProto

	var visit func(protoreflect.FileDescriptor)
	visit = func(current protoreflect.FileDescriptor) {
		if _, ok := seen[current.Path()]; ok {
			return
		}
		seen[current.Path()] = struct{}{}

		imports := current.Imports()
		for i := range imports.Len() {
			visit(imports.Get(i).FileDescriptor)
		}

		descriptors = append(descriptors, protodesc.ToFileDescriptorProto(current))
	}

	visit(file)
	return descriptors
}
