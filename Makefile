#!make

include .env

DB_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable

run:
	go run ./cmd/web

help:
	go run ./cmd/web -h

createDBContainer:
	# docker run --name bookReviewPlatform -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -p 5432:5432 -d postgres
	docker compose up -d --build

createDBPGadmin4Container:
	docker run --name bookPgadmin -p 5050:80 -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin' -d dpage/pgadmin4

createRedisContainer:
	docker run -d --name bookReviewRedis -p 6379:6379 redis:latest 

startContainer:
	docker start bookPgadmin bookReviewRedis

stopContainer:
	docker stop bookPgadmin bookReviewRedis
	
migrateUp: 
	migrate -path migrations -database "${DB_URL}" -verbose up

migrateDown: 
	migrate -path migrations -database "${DB_URL}" -verbose down

migrateForce: 
	migrate -path migrations -database "${DB_URL}" force $(version)

migrateCreate:
	migrate create -ext sql -dir migrations -seq $(fileName)

