#!/bin/bash

set -e

go generate ./...
go run cmd/weasel/main.go
go fmt weasel_otput.go 1>/dev/null 2>/dev/null
cat weasel_otput.go
go run weasel_otput.go
