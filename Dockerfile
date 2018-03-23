FROM golang:1.8
WORKDIR /go/src/github.com/codeuniversity/xing-datahub-producer
RUN curl https://glide.sh/get | sh
COPY . .
RUN make dep
RUN go build producer.go
RUN chmod +x producer
RUN ls
ENTRYPOINT ./producer
EXPOSE 3000
