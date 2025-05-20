BUILD_DIR=./build
BINARY=main
BINARY_PATH=$(BUILD_DIR)/$(BINARY)
SRC=main.go

all: build

build: $(SRC)
	go build -o $(BINARY_PATH) $(SRC)

run: build
	sudo $(BINARY_PATH)

import: 
	@if [ "$(IMG)" = "" ]; then \
		docker save busybox:latest | sudo $(BINARY_PATH); \
	else \
		docker save $(IMG) | sudo $(BINARY_PATH); \
	fi

.PHONY: all build run import
