# Build
FROM golang:1.15-buster AS build

WORKDIR /app
ADD . .

ENV CGO_ENABLED=0

RUN make build

# Enviroment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/bin/mc .
COPY --from=build /app/configs/docker.yml .

ENTRYPOINT ["/app/mc", "-f", "/app/docker.yml"]

EXPOSE 8888
