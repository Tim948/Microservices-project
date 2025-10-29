FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /task-service

EXPOSE 8082

CMD ["/task-service"]