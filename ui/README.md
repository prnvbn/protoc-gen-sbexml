# Proto to SBE XML UI

This is a small browser UI for converting a single pasted `proto3` file into SBE XML. It uses the Go generator through WebAssembly.

## Run

```bash
cd ui
make build
cd dist
python3 -m http.server 8787
```

Open <http://127.0.0.1:8787>.

The UI supports standard protobuf imports provided by the Go protobuf runtime, including `google/protobuf/timestamp.proto`.
