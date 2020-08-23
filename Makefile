SHELL += -eu

BLUE  := \033[0;34m
GREEN := \033[0;32m
RED   := \033[0;31m
NC    := \033[0m

GO111MODULE := on
GOPATH ?= ${HOME}/.gvm/gos/go1.14
GO_BIN := ${GOPATH}/bin
GOPRIVATE := github.com/kevinmichaelchen

# App env vars
KAFKA_HOST ?= kafka

.PHONY: all
all:
	env \
	  GO111MODULE=${GO111MODULE} \
	  KAFKA_HOST=${KAFKA_HOST} \
	  go run main.go

.PHONY: stop
stop:
	docker-compose stop

.PHONY: start
start:
	docker-compose up -d

include makefiles/*.mk