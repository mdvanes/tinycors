# TinyCORS

A tiny CORS Anywhere proxy made with Go.

Adds CORS headers to each request to be able to call APIs that require CORS, without setting up your own server. Just start up this Docker container and start sending requests.

All requests are ... to have ... `Access-Control-Allow-Origin: *!`. Also see [CORS Anywhere](https://github.com/Rob--W/cors-anywhere) and [enable-cors.org](https://enable-cors.org/server.html).

## Why another CORS proxy?

The most popular [CORS Anywhere image](https://hub.docker.com/r/imjacobclark/cors-container/tags) (500k+ downloads at time of writing), uses the `node:10-stretch` image making it 337MB big, which is a lot for what it does.
One of the smaller ones just uses an Nginx configuration (https://hub.docker.com/r/shakyshane/nginx-cors) and is 17MB.

This one uses Go (which compiles to a binary) and is an excellent candidate for [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/). The result is an image of just 12MB!

## Usage

When running on the default port, e.g. go to http://localhost:3000/?url=https://www.mdworld.nl

### Go

`go run tinycors.go`

or with optional port flag:

`go run tinycors.go -port 9000`

or build first:

`go build tinycors.go`
`./tinycors`

### Docker

docker build -t mdworld/tinycors .
docker run --rm --name tinycors -p 3000:3000 mdworld/tinycors

### Docker Compose

## TODO

* add header
* add includelist
* change /?url=x to /x
* pass port flag to docker container
