#
# Hydrowflow
#
HYDROFLOW_BUILDIR?=.build


ifeq ($(OS),Windows_NT)
	uname_S := Windows
else
	uname_S := $(shell uname -s)
endif

.PHONY: app

all: app binary 

clean:
	rm -f .build/*

app: wasm
	python -m http.server 8000 --directory app

binary: $(HYDROFLOW_BUILDIR)/binary

wasm: $(HYDROFLOW_BUILDIR)/hydrowflow.wasm

$(HYDROFLOW_BUILDIR)/binary: clean
	go build -o $(HYDROFLOW_BUILDIR)/hydrowflow cmd/cli/main.go 

$(HYDROFLOW_BUILDIR)/hydrowflow.wasm:
	GOOS=js GOARCH=wasm go build -o $(HYDROFLOW_BUILDIR)/hydrowflow.wasm cmd/wasm/main.go 
