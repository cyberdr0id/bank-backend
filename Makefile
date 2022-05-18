mysql:
	docker run --name mysql-trainee -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql

createdb:
	docker exec -it mysql-trainee mysql -uroot -p -e "create database simple_bank;"

dropdb:
	docker exec -it mysql-trainee mysql -uroot -p -e "USE simple_bank; DROP TABLE IF EXISTS entries; DROP TABLE IF EXISTS transfers; DROP TABLE IF EXISTS accounts;"

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose down

sqlc:
	sqlc generate

.PHONY: mysql createdb dropdb migrateup migratedown sqlc

#
# PostgreSQL
#
# postgres:
# 	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
#
# createdb:
# 	docker exec -i postgres12 createdb --username=root --owner=root simple_bank
#
# dropdb:
#	docker exec -it postgres12 dropdb simple_bank
#
# migrateup:
#	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
#
# migratedown:
#	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
#
