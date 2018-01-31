make dep:
	glide install

make proto:
	protoc --go_out=Protocol Protocol.proto

make run:
	go run main.go

