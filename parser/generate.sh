#!/bin/bash

pigeon ini.peg | goimports | gofmt > pigeon.go
