GO111MODULES=on
APP?=jobcheck
REGISTRY?=docker.io/stevenacoffman
COMMIT_SHA=$(shell git rev-parse --short HEAD)

.PHONY: generate
## generate: generate the application
generate: clean
	@echo "Generating..."
	@go generate ./...

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APP} main.go
	@echo "Done building"

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: docker-build
## docker-build: builds the stringifier docker image to registry
docker-build: build
	docker build -t ${REGISTRY}/${APP}:${COMMIT_SHA} .

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: docker-build
	docker push ${REGISTRY}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
