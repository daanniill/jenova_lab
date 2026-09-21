BINARY := bin/lab
PKG    := ./cmd/lab
ARGS   ?=

.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o $(BINARY) $(PKG)

.PHONY: run
run: build
	./$(BINARY) $(ARGS)
