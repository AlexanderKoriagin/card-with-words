# Commands to build and run cardWithWords application

All commands should be run from the root directory.

## Build docker container

### M1
```shell script
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -mod=vendor -a -o ./docker/cardWithWords ./cmd/cardWithWords
docker buildx build --platform linux/amd64 -o type=docker -t cardwithwords:latest -f ./docker/Dockerfile ./docker
```

### Intel

```shell script
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/cardWithWords ./cmd/cardWithWords
docker build -t cardWithWords:latest -f ./Dockerfile .
```

## Run docker container

```shell script
# Run docker container
docker run -d --name cardWithWords --restart always cardwithwords:latest /cardWithWords --token YOUR_TOKEN
```
