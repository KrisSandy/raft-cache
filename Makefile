all: bin images

build-cli:
	@echo "Building cli..."
	@go build -o bin/raftctl ./cmd/raftctl

build: bin build-cli
	@echo "Building app..."
	@go build -o bin/raft-cache-app ./cmd/app

bin:
	@mkdir -p bin

.PHONY: clean 
clean:
	@rm -r bin/

.PHONY: images
images:
	@echo "Building docker image..."
	@docker build -t krissandy/raft-cache-app .

push-images: images
	@docker push krissandy/raft-cache-app

.PHONY: run
run:
	@go run .

.PHONY: proto
proto:
	protoc ./proto/ca/ca.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative
	protoc ./proto/cs/cs.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative