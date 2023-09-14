#MongoDB container name
DB_CONTAINER_NAME = mongotest

.PHONY: build
build:
	go build -o server cmd/server/*.go

.PHONY: rundb
rundb:
	sudo docker run --name $(DB_CONTAINER_NAME) -d -p 9876:27017 mongodb/mongodb-community-server:latest

.PHONY: stopdb
stopdb:
	sudo docker stop $(DB_CONTAINER_NAME)

.PHONY: killdb
killdb: stopdb
	sudo docker rm $(DB_CONTAINER_NAME)

.PHONY: run
run: build rundb
	./server

.PHONY: reset
reset: killdb run

.PHONY: git
git:
	git add .
	git commit -m "$m"
	git push

.DEFAULT_GOAL := build
