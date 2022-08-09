.PHONY: build

GO:=go
OUTDIR:=bin
LDFLAGS:="-s -w"
ASSETS:=$(shell find -iname '*.go' -type f)

build: $(ASSETS) go.mod go.sum
	$(GO) build \
		-trimpath \
		-ldflags=$(LDFLAGS) \
		-o $(OUTDIR)/clamdiscan \
		cmd/main.go

