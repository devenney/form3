# Payment API

[![Build Status](https://travis-ci.org/devenney/form3.svg?branch=master)](https://travis-ci.org/devenney/form3)

## Configuration

[Viper](https://github.com/spf13/viper) is used for configuration. The important settings, currently set via environment variables, are:

* `PAYMENTS_DB_ENDPOINT`: The DynamoDB endpoint. Leave unset to default to AWS.
* `PAYMENTS_ENV`: The deployment environment. Set to "PROD" to disable dev features.

## Deployment

```shell
make node-deps
make deploy
```

This will deploy the serverless stack with a custom domain (defaults to `f3.devenney.io`). Thus, the endpoints will be available at:

```
GET     https://f3.devenney.io/listPayments
GET     https://f3.devenney.io/getPayment/{id}
POST    https://f3.devenney.io/addPayment
PUT     https://f3.devenney.io/updatePayment/{id}
DELETE  https://f3.devenney.io/deletePayment/{id}
```

## Testing

Testing requires a local instance of DynamoDB, bound to port 8080 (to prevent clashes with the default server port). This is most easily achieved using the official Docker container:

```shell
docker run -p 8080:8000 amazon/dynamodb-local
```

## Local Development

`docker-compose.yml` is provided for easy provisioning of a local development environment.

```shell

### Install docker-compose in Python virtual environment

virtualenv -p python3 ./venv/
. ./venv/bin/activate
pip install docker-compose

### Run docker-compose

docker-compose up
```

## Seeding Data

A command exists to seed DynamoDB, based on the user configuration, with the mock data provided. This expects to run in the root directory of the project as follows:

```shell
go run cmd/seed/main.go
```

If you are attempting to seed a local development deployment, ensure the Docker Compose network is up and running and that the command runs from within the network.

```shell
CONTAINER_ID=$(docker ps -f "name=form3_server" -q)
docker exec -it $CONTAINER_ID /bin/bash
go run cmd/seed/main.go
```
