# [Install golang](https://golang.org/doc/install)

# [Install glide](https://github.com/Masterminds/glide#install)

# Run it
- `make dep`

- `make run`

# Run in Docker
- docker build -t producer .
- docker run -p 3000:3000 --net="host" --rm --name producer producer
