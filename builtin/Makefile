.PHONY: all build
all: build

# $ go tool link
#   ...
#   -s    disable symbol table
#   -w    disable DWARF generation

build:
	go build

.PHONY: run
run: build
	./builtin