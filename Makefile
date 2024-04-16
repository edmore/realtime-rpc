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

create:
		docker build -f Dockerfile.terraform -t pennsieve/realtime-rpc-deploy . --progress=plain
		docker-compose run realtime-rpc-deploy create

status:
		docker build -f Dockerfile.terraform -t pennsieve/realtime-rpc-deploy . --progress=plain
		docker-compose run realtime-rpc-deploy status


destroy:
		docker build -f Dockerfile.terraform -t pennsieve/realtime-rpc-deploy . --progress=plain
		docker-compose run realtime-rpc-deploy destroy


deploy:
		aws ecr get-login-password --profile ${AWS_PROFILE} --region ${AWS_DEFAULT_REGION} | docker login --username AWS --password-stdin ${ACCOUNT}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com
		docker buildx build --platform linux/amd64 --progress=plain -t pennsieve/realtime-rpc .
		docker tag pennsieve/realtime-rpc ${APP_REPO}
		docker push ${APP_REPO}

