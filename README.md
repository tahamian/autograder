# Autograder

## Run using Docker compose

```
cd autograder


```

Requires a redis instance to running set in config.yaml

## Run as a docker contianer

```
cd autograder/server/

docker build . -t autograder

```


## Build server

```
go get -d autograder

cd autograder/server/ 

make all
```

## Run go server

```
go get -d autograder

go run main.go
```