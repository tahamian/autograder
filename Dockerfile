FROM golang:alpine as base

MAINTAINER MostafaAyesh

COPY ./* /autograder/

# install dependencies
RUN    apk add --update \
    && apk upgrade \
    && apk add --no-cache bash git make

# build autograder
RUN    cd /autograder \ 
    && make \
    && cp -r bin/* /usr/bin

ENTRYPOINT ["/usr/bin/autograder", "--config=/autograder/config.yaml"]
