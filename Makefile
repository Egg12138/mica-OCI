BINARY=main
BINARY_PATH=./$(BINARY)
SRC=main.go

all: build

build: $(SRC)
	go build -o $(BINARY) $(SRC)

run: build
	sudo $(BINARY_PATH)

import: 
	@if [ "$(IMG)" = "" ]; then \
		docker save busybox:latest | sudo $(BINARY_PATH); \
	else \
		docker save $(IMG) | sudo $(BINARY_PATH); \
	fi

.PHONY: all build run import