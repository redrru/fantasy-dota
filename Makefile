all: run

test:
	go test -tags unit ./...

style:
	go fmt ./...
	go mod tidy
	golangci-lint run

run:
	docker-compose -f ./build/docker/docker-compose.yaml up --build

down:
	docker-compose -f ./build/docker/docker-compose.yaml down

codegen:
	oapi-codegen -old-config-style -generate "types,server" -package server api/http/openapi.yaml > pkg/server/fantasy-dota.gen.go
