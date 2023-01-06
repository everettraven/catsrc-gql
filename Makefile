SHELL = /bin/bash

.PHONY: build
build:
	go build -o catsrc-gql main.go