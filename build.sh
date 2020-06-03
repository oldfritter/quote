#!/bin/sh
go build -o ./cmd/getData sourceData/getData.go
go build -o ./cmd/api api/api.go
go build -o ./cmd/worker workers/worker.go
go build -o ./cmd/schedule schedules/schedule.go
