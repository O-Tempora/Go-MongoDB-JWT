#MongoDB container name
DB_CONTAINER_NAME = mongotest

build:
	go build -o server cmd/server/*.go

rundb:
	sudo docker run --name $(DB_CONTAINER_NAME) -d -p 9876:27017 mongodb/mongodb-community-server:latest
stopdb:
	sudo docker stop $(DB_CONTAINER_NAME)
killdb: stopdb
	sudo docker rm $(DB_CONTAINER_NAME)

run: build rundb
	./server
