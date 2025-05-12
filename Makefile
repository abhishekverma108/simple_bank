postgresup:
	docker compose -f /home/roadcast/go_projects/tech-school_go_database/docker-compose.yml  up -d 
postgresdown:
	docker compose -f /home/roadcast/go_projects/tech-school_go_database/docker-compose.yml  down
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: postgresup postgresdown createdb dropdb migratedown migrateups sqlc



