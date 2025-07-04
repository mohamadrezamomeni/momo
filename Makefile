OS := $(shell uname -s)
EXT :=
ifeq ($(OS),Windows_NT)
	EXT := .exe
endif

CMD_DIR := ./cmd
BIN_DIR := ./bin

BINARIES := $(shell find $(CMD_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

TARGETS := $(foreach bin,$(BINARIES),$(BIN_DIR)/$(bin)$(EXT))

.PHONY: all build clean

all: build

build: $(TARGETS)

$(BIN_DIR)/%$(EXT):
	@echo "Building $(@F)..."
	go build -o $@ $(CMD_DIR)/$*

clean:
	rm -rf $(BIN_DIR)
