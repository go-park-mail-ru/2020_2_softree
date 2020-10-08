# Build
FROM golang:1.15-buster AS build

WORKDIR /app
ADD . .

RUN make build

# Enviroment
FROM debian:buster-slim

WORKDIR /app
COPY --from=build /app/bin/mc .
COPY --from=build /app/docker.yml .

ENTRYPOINT ["/app/mc", "-f", "app/docker.yml"]

EXPOSE 8888
