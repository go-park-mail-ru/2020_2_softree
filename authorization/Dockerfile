# Build
FROM golang:1.15-buster AS build

WORKDIR /app
ADD . .

ENV CGO_ENABLED=0

RUN make TARGET="authorization" build

# Enviroment
FROM alpine:latest

WORKDIR /app
RUN mkdir /etc/moneycat
COPY --from=build /app/bin/authorization .

ENTRYPOINT ["/app/authorization"]

EXPOSE 8000