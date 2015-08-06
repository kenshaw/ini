#!/bin/bash

go test -v -covermode=count -coverprofile=coverage.out -coverpkg=github.com/knq/ini/parser && go tool cover -html=coverage.out
