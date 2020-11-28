#!/bin/bash

go test -count=1 -coverpkg=./... -coverprofile=cover -v -race -timeout 30s ./...
# shellcheck disable=SC2002
# shellcheck disable=SC2197
cat cover | fgrep -v "mock" > cover_wo_mock
go tool cover -func cover_wo_mock
rm cover cover_wo_mock
