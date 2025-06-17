build-debug:
	@go build -gcflags "all=-N -l" -o tmp/main .

local-debug: build-debug
	@dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./tmp/main

build:
	@go build -o tmp/main .

local: build
	@./tmp/main

install:
	@go install .
