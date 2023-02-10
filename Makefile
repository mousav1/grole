postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root grole
	
test:
	go test -v -cover ./...


.PHONY: postgres createdb test
