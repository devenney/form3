FROM golang

WORKDIR /go/src/app

ENV GO111MODULE=on

# TODO(devenney): Hot code reloading. Would be nicer to mount than copy at build time.
COPY . .

RUN go install -v ./...

CMD server
