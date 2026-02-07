
.PHONY: createdb dropdb migrate-up migrate-down sqlc test

createdb:
	docker-compose up -d postgres
	sleep 2
	docker-compose exec postgres createdb -U admin simple_bank

dropdb:
	docker-compose down -v
	docker-compose exec postgres dropdb -U admin simle_banck --if-exists --force

migrate-up:
	migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" down

migrate-force:
	migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" force $(version)

migrate-reset:
	migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" drop -f
	migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" up


sqlc:
	sqlc generate

test:
	go test -v -cover ./...