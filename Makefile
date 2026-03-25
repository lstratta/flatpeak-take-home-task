# test the application code
# pgadmin does not need to run
# cleanup first to remove any leftover data
test: cleanup
	@go test -v ./...
.PHONY: test

# build the binary after testing the code
build: test
	@go build -o main cmd/main.go
.PHONY: build

# simply starts the database containers and the development server 
start: 
	@air
.PHONY: start

# starts the server without Air
run:
	@go run cmd/main.go
.PHONY: run 

# stops and removes all containers 
cleanup: 
	@echo "stopping container..."
	@docker stop carbon-intensity -i
	@echo "removing container..."
	@docker rm pgadmin carbon-intensity -i
.PHONY: cleanup

docker-build: test
	@docker build . -t flatpeak/carbon-intensity:latest

docker-run: docker-build
	@docker compose -f docker/compose.yaml up -d 
.PHONY: docker-run
