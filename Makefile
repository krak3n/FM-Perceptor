#
# Makefile for FM-Perceptor
#
# Handy shortcuts for running go tool commands
#

build:
	go build ./cmd/perceptor

install:
	go install ./cmd/perceptor
