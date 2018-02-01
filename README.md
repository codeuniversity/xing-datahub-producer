# [Install golang](https://golang.org/doc/install)

# [Install glide](https://github.com/Masterminds/glide#install)

# Run it
- `make dep`

- `make run`

# Run in Docker
- In one terminal: `docker-compose up`
- In another terminal: `docker-compose exec go bash`
- Now you can use `make dep` and `make run` to install dependencies and run
  your code.
- Note: when running `make dep` you may get asked to login to your github
  account to get access to the codeuniversity repository, do so using as
  password the token generated for single log on.
