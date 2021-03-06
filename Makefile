# mysql:
# 	docker run --name mysql-trainee -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql

# createdb:
# 	docker exec -it mysql-trainee mysql -uroot -p -e "create database simple_bank;"

# dropdb:
# 	docker exec -it mysql-trainee mysql -uroot -p -e "USE simple_bank; DROP TABLE IF EXISTS entries; DROP TABLE IF EXISTS transfers; DROP TABLE IF EXISTS accounts;"

# migrateup:
# 	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose up

# migratedown:
# 	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose down


#
# PostgreSQL
#
postgres:
	docker run --name postgres-bank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -i postgres-bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-bank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test.unit:
	go test -v -cover ./db/sqlc/...

test.api:	
	go test -v -cover ./api/...
	
server:
	go run cmd/main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/cyberdr0id/bank-backend/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown sqlc server mock
