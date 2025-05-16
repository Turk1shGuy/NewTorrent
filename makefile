GO_BUILD_CMD = go build
GO_RUN_CMD = go run
BINARY_NAME = torrent_server.elf
BINARY_PATH = ./cmd/torrent
CRAWL_PATH	= ./crawler/crawl.go

build:
	$(GO_BUILD_CMD) -o $(BINARY_NAME) $(BINARY_PATH)

run:
	$(GO_RUN_CMD) $(BINARY_PATH)

runc:
	$(GO_RUN_CMD) $(CRAWL_PATH)

clean:
	rm -f $(BINARY_NAME)

all: build run

.PHONY: build run clean all
