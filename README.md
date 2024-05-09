A POC to run an RPC server and client on AWS:

Based on: https://aws.amazon.com/blogs/opensource/containerize-and-deploy-a-grpc-application-on-aws-fargate/

**To run locally:**

*Build:*

`docker build -t jit_server .`

`docker build --no-cache -t jit_client client/.`

`docker run -d --name jit_server jit_server`

`docker run -itd --name jit_client --link jit_server --env "SERVER_ENDPOINT=jit_server:50051" jit_client`


 *To view output:*
 ```
 docker logs jit_server
 docker logs jit_client
 ```

 You can also opt to run the images on separate terminals:
 `docker run --name jit_server jit_server`
 `docker run -it --name jit_client --link jit_server --env "SERVER_ENDPOINT=jit_server:50051" jit_client`

 **To run on the cloud:**

 1. Rename `infa.env.sample` to `infra.dev`
 2. Pre-requisites: AWS (prerequisite(s): [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html))
 3. Update the env file with your AWS [profile](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-methods) and region
 4. To create the infrastructure: `make create`
 5. To deploy:

 ```
 make deploy ACCOUNT=<account> AWS_DEFAULT_REGION=<region> APP_REPO=<account>.dkr.ecr.<region>.amazonaws.com/realtime-rpc AWS_PROFILE=<profile> CLIENT_REPO=<account>.dkr.ecr.<region>.amazonaws.com/realtime-rpc-client
 ```

 6. Once deployed. Navigate to Amazon Elastic Service > Task Definitions > realtime-rpc-client > select and "run task" > select cluster and "create" (this will invoke the task) > check client and server logs for result.


 **Future Work**

 Connecting directly to the RPC server from the internet: https://aws.amazon.com/blogs/aws/new-application-load-balancer-support-for-end-to-end-http-2-and-grpc/