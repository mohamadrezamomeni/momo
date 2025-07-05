CMD_DIR := ./cmd
BIN_DIR := ./bin

BINARIES := $(shell find $(CMD_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
TARGETS := $(foreach bin,$(BINARIES),$(BIN_DIR)/$(bin))

.PHONY: all build clean

all: build

build: $(TARGETS)

$(BIN_DIR)/%:
	@echo "🔨 Building $(@F)..."
	@mkdir -p $(BIN_DIR)
	go build -o $@ $(CMD_DIR)/$*

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(BIN_DIR)