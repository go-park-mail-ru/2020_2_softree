#!/bin/bash

go test -count=1 -coverprofile=cover -v -race -timeout 30s ./...

grep -F -v "mock" cover > cover_wo_mock
go tool cover -func cover_wo_mock

rm cover cover_wo_mock
