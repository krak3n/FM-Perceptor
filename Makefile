#
# Makefile for FM-Perceptor
#
# Handy shortcuts for running go tool commands
#

build:
	go build $(glide nv)

install:
	go install $(glide nv)
