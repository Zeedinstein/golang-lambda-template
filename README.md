# Golang lambda template

## Run DynamoDB
Run docker-compose to start up DynamoDB
```
docker-compose up
```
### Build Lambda
Build code and Cloudformation template
```
make build
```
### Run Lambda
Run code locally
```
sam local start-api --docker-network YOUR_DOCKER_NETWORK
```
