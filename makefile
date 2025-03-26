SHELL := /bin/bash

.PHONY: build up down logs psql migrate run

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

psql:
	docker-compose exec db psql -U postgres -d clicks

migrate:
	docker-compose exec db psql -U postgres -d clicks -f migrations/001_init.sql

run:
	go run cmd/server/main.go
