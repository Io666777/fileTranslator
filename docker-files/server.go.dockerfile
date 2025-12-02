FROM golang:1.25.3

RUN apt-get update && apt-get install -y tzdata

WORKDIR /translater_server

ARG SERVER_MODULE_PATH
ARG SERVER_SERVICE_PATH

COPY ./project/ ./

COPY .env ./

RUN go mod download

RUN go build -o main ./cmd

CMD ["/translater_server/main"]