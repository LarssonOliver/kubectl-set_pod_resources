# Copyright 2025 Oliver Larsson

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and

BIN_DIR := bin
BIN_NAME := kubectl-set_pod_resources

SOURCE_FILES := go.mod go.sum $(shell find cmd pkg -name '*.go')

all: $(BIN_DIR)/$(BIN_NAME)

$(BIN_DIR)/$(BIN_NAME): $(SOURCE_FILES)
	go build -o $@ cmd/kubectl-set_pod_resources.go

.PHONY: clean

clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)
