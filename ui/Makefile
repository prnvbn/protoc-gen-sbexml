.PHONY: build clean serve copy-vendor

GO ?= go
DIST_DIR := dist

build: $(DIST_DIR)/index.html $(DIST_DIR)/app.js $(DIST_DIR)/loader.js $(DIST_DIR)/styles.css copy-vendor $(DIST_DIR)/version.txt $(DIST_DIR)/wasm_exec.js $(DIST_DIR)/main.wasm

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

$(DIST_DIR)/index.html: index.html | $(DIST_DIR)
	cp index.html $(DIST_DIR)/index.html

$(DIST_DIR)/app.js: app.js | $(DIST_DIR)
	cp app.js $(DIST_DIR)/app.js

$(DIST_DIR)/loader.js: loader.js | $(DIST_DIR)
	cp loader.js $(DIST_DIR)/loader.js

$(DIST_DIR)/styles.css: styles.css | $(DIST_DIR)
	cp styles.css $(DIST_DIR)/styles.css

$(DIST_DIR)/version.txt: ../version.txt | $(DIST_DIR)
	cp ../version.txt $(DIST_DIR)/version.txt

copy-vendor: | $(DIST_DIR)
	rm -rf $(DIST_DIR)/vendor
	cp -R vendor $(DIST_DIR)/vendor

$(DIST_DIR)/wasm_exec.js: | $(DIST_DIR)
	cp "$$($(GO) env GOROOT)/lib/wasm/wasm_exec.js" $(DIST_DIR)/wasm_exec.js

$(DIST_DIR)/main.wasm: main.go ../internal/protosource/generate.go ../internal/sbexml/*.go | $(DIST_DIR)
	GOOS=js GOARCH=wasm $(GO) build -o $(DIST_DIR)/main.wasm .

serve: build
	cd $(DIST_DIR) && python3 -m http.server 8787

clean:
	rm -rf $(DIST_DIR)
