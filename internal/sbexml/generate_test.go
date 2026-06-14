package sbexml_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prnvbn/protoc-gen-sbexml/internal/sbexml"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerate(t *testing.T) {
	for _, tc := range discoverTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			request := compileProtoCase(t, tc.inputDir)

			got, err := sbexml.Generate(request)
			require.NoError(t, err)

			wantBytes, err := os.ReadFile(tc.outputFile)
			require.NoError(t, err)
			want := string(wantBytes)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatalf("generated XML mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func discoverTestCases(t *testing.T) []testCase {
	t.Helper()

	entries, err := os.ReadDir("testcases")
	require.NoError(t, err)

	var cases []testCase
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join("testcases", entry.Name())
		cases = append(cases, testCase{
			name:       entry.Name(),
			inputDir:   filepath.Join(dir, "input"),
			outputFile: filepath.Join(dir, "output.xml"),
		})
	}

	return cases
}

type testCase struct {
	name       string
	inputDir   string
	outputFile string
}

func compileProtoCase(t *testing.T, inputDir string) *pluginpb.CodeGeneratorRequest {
	t.Helper()

	protoFiles := protoFilesToGenerateIn(t, inputDir)
	descriptorPath := filepath.Join(t.TempDir(), "descriptor.pb")

	args := []string{
		"--proto_path=" + inputDir,
		"--include_imports",
		"--descriptor_set_out=" + descriptorPath,
	}
	args = append(args, protoFiles...)

	cmd := exec.Command("protoc", args...)
	output, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "protoc failed:\n%s", output)

	descriptorBytes, err := os.ReadFile(descriptorPath)
	require.NoError(t, err)

	var descriptorSet descriptorpb.FileDescriptorSet
	require.NoError(t, proto.Unmarshal(descriptorBytes, &descriptorSet))

	return &pluginpb.CodeGeneratorRequest{
		FileToGenerate: protoFiles,
		ProtoFile:      descriptorSet.File,
	}
}

func protoFilesToGenerateIn(t *testing.T, inputDir string) []string {
	t.Helper()

	var protoFiles []string
	entries, err := os.ReadDir(inputDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".proto" {
			continue
		}
		protoFiles = append(protoFiles, entry.Name())
	}

	require.NotEmpty(t, protoFiles, "no root .proto files found in %s", inputDir)
	slices.Sort(protoFiles)
	return protoFiles
}
