# Autograder

[![Build Status](https://travis-ci.com/tahamian/autograder.svg?branch=master)](https://travis-ci.com/tahamian/autograder)
[![Go Report Card](https://goreportcard.com/badge/github.com/tahamian/autograder)](https://goreportcard.com/report/github.com/tahamian/autograder)

Autograder evalutes python scripts uploaded to the server and checks the output.
All evaluations are done in docker containers for code isolation.
If a user submits malisious code it will be excecuted inside the container.


## Requirements to run

- Docker installed
- npm installed

## Run app

### Run using Docker compose


```
cd autograder


```

Requires a redis instance to running set in config.yaml

### Run as a docker contianer

```
cd autograder/server/

docker build . -t autograder

```


#### Build server

```
go get -d autograder

cd autograder/server/ 

make all
```

#### Run go server

```
go get -d autograder

go run main.go
```

##### Installing js libraries

```$xslt
These are js libraries for the ui

# install gulp
npm install gulp --save-dev


# install bower
npm install -g bower

# run in project dir
gulp
```

## Run Tests

```
cd autograder

go test
```
