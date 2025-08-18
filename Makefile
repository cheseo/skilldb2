.PHONY: run db

run: skilldb2
	./skilldb2

skilldb2: *.go
	go build

db:
	-rm db.db
	./createdb
	-sqlite3 db.db < sample.sql
