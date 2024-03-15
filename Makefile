.DEFAULT: help

SERVICE_NAME  ?= "jit-calculation-service"

help:
	@echo "Make Help for $(SERVICE_NAME)"

compile:
		docker run -v `pwd`:/tmp edmore/protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative \
   		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	--proto_path=.

test:
		go test -race ./...