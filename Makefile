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
	cd web && npm run build

.PHONY: clear
clear: 
	rm -r web/build && rm -r webhook-listener && rm -r web.tar.gz

.PHONY: pack
pack:
	cd web && npm install && npm run build:tailwind && REACT_APP_API_PROTOCOL=https REACT_APP_API_HOST=webhook.monster npm run build && \
	cd ../ && tar -czf web.tar.gz ./web/build && \
	CGO_ENABLED=0 GOOS=linux go build \