DDB_PORT=8080
DDB_CONTAINER_NAME=form3_payments_testing_ddb

ci: clean lint test

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/list lambda/list/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get lambda/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add lambda/add/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update lambda/update/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete lambda/delete/main.go

deploy: build
	sls create_domain
	sls deploy

node-deps:
	npm install --dev serverless
	npm install --dev serverless-domain-manager

install:
	go install -v ./...

lint:
	$(GOBIN)/golint -set_exit_status ./...

test: clean
	docker run -p $(DDB_PORT):8000 -d --name $(DDB_CONTAINER_NAME) amazon/dynamodb-local
	go test -cover -v ./...


clean:
	rm -rf node_modules
	-docker stop $(DDB_CONTAINER_NAME)
	-docker rm $(DDB_CONTAINER_NAME)
