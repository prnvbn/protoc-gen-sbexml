.PHONY: dist windows darwin-amd64 darwin-arm64 linux-arm64 linux-amd64 clean

DIST_DIR := dist
LDFLAGS := -s -w
BUILD_FLAGS := -ldflags="$(LDFLAGS)"
CMD := ./cmd/protoc-gen-sbexml
BINARY := protoc-gen-sbexml

define build-target
CGO_ENABLED=0 GOOS=$(1)$(if $(2), GOARCH=$(2)) \
	go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY)-$(3) $(CMD)
endef

dist:
	mkdir -p $(DIST_DIR)
	$(MAKE) windows darwin-arm64 darwin-amd64 linux-arm64 linux-amd64

windows:
	$(call build-target,windows,amd64,windows-amd64.exe)

darwin-amd64:
	$(call build-target,darwin,amd64,darwin-amd64)

darwin-arm64:
	$(call build-target,darwin,arm64,darwin-arm64)

linux-arm64:
	$(call build-target,linux,arm64,linux-arm64)

linux-amd64:
	$(call build-target,linux,amd64,linux-amd64)

clean:
	rm -rf $(DIST_DIR)
