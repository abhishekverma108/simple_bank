postgresup:
	docker compose -f /home/roadcast/go_projects/tech-school_go_database/docker-compose.yml  up -d 
postgresdown:
	docker compose -f /home/roadcast/go_projects/tech-school_go_database/docker-compose.yml  down
createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://postgres:pae9bai7Cahg?ahcae\"g@134.209.150.195:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres:pae9bai7Cahg?ahcae\"g@134.209.150.195:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: postgresup postgresdown createdb dropdb migratedown migrateups sqlc server



