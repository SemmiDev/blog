build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/main main.go
run-on-container:
	docker build -t app:v1 . && \
	docker run -d --name web1 --net blog_mynetwork app:v1 && \
	docker run -d --name web2 --net blog_mynetwork app:v1 && \
	docker run -d --name web3 --net blog_mynetwork app:v1
run:
	go run main.go

.PHONY: build run run-on-container