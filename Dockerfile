# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base gcc
ADD . /src
RUN cd /src && go build -o tinycors

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/tinycors /app/
ENTRYPOINT ["./tinycors"]