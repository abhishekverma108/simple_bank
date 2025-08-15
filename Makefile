postgresup:
	docker compose -f /Users/roadcast/roadcast_old_laptop/go_projects/tech-school_go_database/docker-compose.yml  up -d
postgresdown:
	docker compose -f /Users/roadcast/roadcast_old_laptop/go_projects/tech-school_go_database/docker-compose.yml  down
createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5445/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5445/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb  -destination db/mock/store.go simplebank/db/sqlc Store


.PHONY: postgresup postgresdown createdb dropdb migratedown migrateups sqlc server



