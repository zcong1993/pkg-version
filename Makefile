generate:
	@go generate ./...

build: generate
	@echo "====> Build pkg-version"
	@sh -c ./build.sh
