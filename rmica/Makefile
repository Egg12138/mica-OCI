BINDIR ?= ./build
BINNAME ?= rmica
BINARY_PATH := $(BINDIR)/$(BINNAME)
VERSION ?= $(shell cat VERSION)
GO ?= go
DCKR := docker -D
LDFLAGS ?= "-w -s"

# M=mode of test/delete target
# R=value of --runtime=
# C=container name you wanna delete
# K=keep the container created (default to be false)

R ?= rmica
SHELL_IMG ?= busybox
SIMPLE_IMG ?= hello-world

# IMG defaults to SIMPLE_IMG, but can be overridden by M=it
IMG := $(SIMPLE_IMG)
ifeq ($(M),it)
IMG := $(SHELL_IMG)
endif

# RUN_OPTS depends on R
RUN_OPTS ?= --runtime=$(R)
ifeq ($(M),it)
RUN_OPTS := $(RUN_OPTS) -it
endif

ifeq ($(K),)
RUN_OPTS := $(RUN_OPTS) --rm
endif

# CNTR_NAME depends on final IMG value
CNTR_NAME := $(IMG)_container

.PHONY: all build run clean delete delete_all test help

all: build

help:
	@echo "Usage: make [target] [R=<runtime>] [M="it"] [C=<container>]"
	@echo ""
	@echo "Targets:"
	@echo "  all        Build binary (default)"
	@echo "  build      Build binary"
	@echo "  run        Run binary"
	@echo "  clean      Remove build artifacts"
	@echo "  test       Test container run"
	@echo "  delete     Delete container"
	@echo "  delete_all Delete all containers"
	@echo ""
	@echo "Variables:"
	@echo "  R          Runtime to use (default: rmica)"
	@echo "  M          Mode: 'it' for interactive"
	@echo "  C          Container name to delete"

build:
	mkdir -p $(BINDIR)
	$(GO) build -o $(BINARY_PATH)

$(BINNAME): build
	$(BINARY_PATH)

clean:
	rm -f $(BINARY_PATH)

run: build
	$(DCKR) run $(RUN_OPTS) --name $(CNTR_NAME) $(IMG)

delete_all:
	@echo "removing $(IMG) container instances"
	@echo '$(DCKR) rm $$(docker ps -a --filter "ancestor=$(IMG):latest" --filter "status=created" -q)'

delete:
ifneq ($(C),)
	$(DCKR) rm $(C)
else
	$(DCKR) rm $(CNTR_NAME)
endif
