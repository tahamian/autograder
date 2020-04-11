# Autograder

[![Build Status](https://travis-ci.com/tahamian/autograder.svg?branch=master)](https://travis-ci.com/tahamian/autograder)


Autograder evalutes python scripts uploaded to the server and checks the output.
All evaluations are done in docker containers for code isolation.
If a user submits malicious code it will be executed inside the container.


## Requirements to run

- Docker installed
- npm installed

## Run app

*If you want to run a redis instance set in config.yaml*

### Run as a docker contianer

```
cd autograder/server/

docker build . -t autograder
```


#### Build server

```
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

## How Configure tests

In `config.yaml` 


Example config of a program testing stdout
```$xslt
labs:
    - 
```

Example config of a program testing functions
```$xslt
labs:
    -
```


Multiple expected values
```$xslt
labs:
    -
```

## Run Tests

### Run go server test
```
go test
```

### Run python marker tests
```$xslt
python -m pytest marker/tests
```



