# protoc-gen-sbexml

`protoc-gen-sbexml` is a `protoc` plugin that generates [Simple Binary Encoding](https://github.com/aeron-io/simple-binary-encoding) XML from protobuf descriptors.

## Installation

### via `homebrew`

You can install `protoc-gen-sbexml` using the [prnvbn/homebrew-tap](https://github.com/prnvbn/homebrew-tap).

```bash
brew install prnvbn/tap/protoc-gen-sbexml
```

### via `go install`

```bash
go install github.com/prnvbn/protoc-gen-sbexml/cmd/protoc-gen-sbexml@latest
```

### From source

```bash
git clone https://github.com/prnvbn/protoc-gen-sbexml.git
cd protoc-gen-sbexml
go build -o /usr/local/bin/protoc-gen-sbexml ./cmd/protoc-gen-sbexml
```

`protoc` discovers plugins by name, so make sure the installed `protoc-gen-sbexml` binary is available in your `PATH`.

## WIP Status

> [!WARNING]
>
> This generator is a work in progress. The XML output is intentionally minimal while the proto-to-SBE mapping is still being designed.

- [x] proto3 descriptor input through `protoc`
- [x] packages, top-level messages, and nested messages
- [x] top-level enums, nested enums, and enum fields
- [x] scalar primitive fields
- [x] repeated primitive, enum, and message fields
- [x] oneof fields
- [ ] imports and cross-file type references
- [x] maps
- [ ] proto3 `optional`
- [ ] reserved fields
- [ ] well known types
- [ ] additional proto annotations (to support better memory efficieny)
