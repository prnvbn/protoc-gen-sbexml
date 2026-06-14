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
