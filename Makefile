SOURCEDIR=cmd
BINARY=ptpip
BINARY_NOLV=${BINARY}-nolv
VERSION := $(shell git describe --tags)
BUILD_TIME := $(shell date +%FT%T%z)

LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
TAGS=-tags with_lv

.DEFAULT_GOAL: all

.PHONY: all
all: ptpip

ptpip:
	cd cmd; go build ${LDFLAGS} ${TAGS} -o ../${BINARY}

nolv:
	cd cmd; go build ${LDFLAGS} -o ../${BINARY_NOLV}

.PHONY: install
install:
	cd cmd; GOBIN=/usr/local/bin/ go install ${LDFLAGS} ${TAGS}

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi ; if [ -f ${BINARY_NOLV} ] ; then rm ${BINARY_NOLV} ; fi
