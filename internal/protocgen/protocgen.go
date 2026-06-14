package protocgen

import (
	"fmt"
	"io"

	"github.com/prnvbn/protoc-gen-sbexml/internal/sbexml"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// Run reads a protoc plugin request from in and writes a plugin response to out.
func Run(in io.Reader, out io.Writer) error {
	response, err := Generate(in)
	if err != nil {
		response = &pluginpb.CodeGeneratorResponse{
			Error: new(err.Error()),
		}
	}

	data, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	if _, err := out.Write(data); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

// Generate converts a serialized protoc plugin request into a plugin response.
func Generate(in io.Reader) (*pluginpb.CodeGeneratorResponse, error) {
	request, err := readRequest(in)
	if err != nil {
		return nil, fmt.Errorf("read code generator request: %w", err)
	}

	content, err := sbexml.Generate(request)
	if err != nil {
		return nil, fmt.Errorf("generate SBE XML: %w", err)
	}

	supportedFeatures := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	return &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &supportedFeatures,
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    new(sbexml.OutputFilename),
				Content: new(content),
			},
		},
	}, nil
}

func readRequest(in io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("read request: %w", err)
	}

	var request pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("unmarshal request: %w", err)
	}
	return &request, nil
}
