runpostgres:
	docker run --name postgre1 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -p 5432:5432 -d postgres:14.5-alpine

createdb:
	docker exec -it postgre1 createdb --user=root --owner=root simplebank

dropdb:
	docker exec -it postgre1 dropdb simplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: runpostgres createdb dropdb migrateup migratedown sqlc