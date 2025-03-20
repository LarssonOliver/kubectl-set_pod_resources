BIN_DIR := bin
BIN_NAME := kubectl-resize_pod

SOURCE_FILES := $(shell find cmd pkg -name '*.go')

all: $(BIN_DIR)/$(BIN_NAME)

$(BIN_DIR)/$(BIN_NAME): $(SOURCE_FILES)
	go build -o $@ cmd/kubectl-resize_pod.go

.PHONY: clean

clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)
