SOURCEDIR=cmd
SOURCES := ${SOURCEDIR}/*.go
BINARY=ptpip
VERSION := $(shell git describe --tags)
BUILD_TIME := $(shell date +%FT%T%z)

LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: all

.PHONY: all
all: ptpip

ptpip:
	go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	go install ${LDFLAGS} ${SOURCES}

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
