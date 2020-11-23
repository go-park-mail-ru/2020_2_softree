#!/bin/bash

go test -count=1 -coverpkg=./... -coverprofile=cover -v -race -timeout 30s ./...
cat cover | fgrep -v "mock" > cover_wo_mock
go tool cover -func cover_wo_mock
rm cover cover_wo_mock
