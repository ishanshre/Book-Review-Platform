#!make

include .env

DB_URL=postgresql://${m_db_username}:${m_db_password}@localhost:5432/${m_db_dbname}?sslmode=disable

run:
	go run ./cmd/web 8000

createDBContainer:
	docker run --name bookReviewPlatform -e POSTGRES_USER=${m_db_username} -e POSTGRES_PASSWORD=${m_db_password} -p 5432:5432 -d postgres

createDBPGadmin4Container:
	docker run --name pgadmin -p 5050:80 -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin' -d dpage/pgadmin4

createRedisContainer:
	docker run -d --name bookReviewRedis -p 6379:6379 redis:latest 

migrateUp: 
	migrate -path migrations -database "${DB_URL}" -verbose up

migrateDown: 
	migrate -path migrations -database "${DB_URL}" -verbose down

migrateForce: 
	migrate -path migrations -database "${DB_URL}" force $(version)

migrateCreate:
	migrate create -ext sql -dir migrations -seq $(fileName)

