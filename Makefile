all: run

test:
	go test -tags unit ./...

style:
	go fmt ./...
	go mod tidy
	golangci-lint run

up:
	docker-compose -f ./build/docker/docker-compose.yaml up --build

down:
	docker-compose -f ./build/docker/docker-compose.yaml down

codegen:
	oapi-codegen -old-config-style -generate "types,server" -package server api/http/openapi.yaml > pkg/server/fantasy-dota.gen.go

run-compile-daemon:
	CompileDaemon \
		-build="go build -o ./.tmp/fantasy-dota ./cmd/fantasy-dota/main.go" \
		-command="./.tmp/fantasy-dota" \
		-exclude-dir=api \
		-exclude-dir=.run \
		-exclude-dir=.git \
		-exclude=fantasy-dota \
		-color \
		-graceful-kill=true
