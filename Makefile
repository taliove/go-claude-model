install:
	go build -ldflags="-s -w" -o ccm .
	cp ccm /home/tangfan/claude-model/bin/ccm

build:
	go build -ldflags="-s -w" -o ccm .

.PHONY: install build
