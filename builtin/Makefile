.PHONY: all build
all: build

# $ go tool link
#   ...
#   -s    disable symbol table
#   -w    disable DWARF generation

build:
	CGO_ENABLED=0 go build -ldflags="-s -w"
	# upx --best --lzma builtin