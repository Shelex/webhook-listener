NAME=webhook-listener
ROOT=github.com/Shelex/${NAME}
GO111MODULE=on
SHELL=/bin/bash

.PHONY: start
start:
	go run main.go

.PHONY: simulation
simulation:
	k6 run simulation.js

.PHONY: open
open-web:
	open http://localhost:8080

.PHONY: prof
prof:
	go tool pprof -web http://localhost:6060/debug/pprof/heap

.PHONY: doc
doc:
	swag init

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: web-dev
web-dev: 
	cd web && npm start

.PHONY: web-build
web-build: 
	cd web && npm build