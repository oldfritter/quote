#!/bin/sh
go build -o ./cmd/getData sourceData/getData.go
go build -o ./cmd/api api/api.go
