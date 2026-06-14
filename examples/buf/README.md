# Buf Example

This example runs an installed `protoc-gen-sbexml` binary as a local [Buf](https://github.com/bufbuild/buf) plugin.

Make sure `protoc-gen-sbexml` is available in your `PATH`.

From this directory:

```bash
buf generate
```

Buf uses `buf.gen.yaml` to run the plugin with:

```yaml
local: protoc-gen-sbexml
```

The generated SBE XML is written to:

```text
gen/output.xml
```
