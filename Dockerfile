FROM golang:1.8
WORKDIR /go/src/github.com/codeuniversity/xing-datahub-producer
RUN curl https://glide.sh/get | sh
EXPOSE 8080
