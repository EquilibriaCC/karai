FROM        golang:alpine AS builder

MAINTAINER  rocksteadytc rock@karai.io

WORKDIR     /home/karai
ADD         . /home/karai

RUN         apk add git              && \
            go build
    
FROM        alpine

COPY        --from=builder              \
            /home/harrisonhesslink/pythia go-karai

EXPOSE      4200

ENTRYPOINT  ["./go-karai"]