FROM golang:1.8
WORKDIR /go/src/xing-datahub-producer
RUN curl https://glide.sh/get | sh
