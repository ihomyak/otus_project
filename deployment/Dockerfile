FROM golang:1.22.8-alpine3.19 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/server /app/cmd/main.go

FROM alpine:3.19

ARG SERVER_PORT

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env.example /app/.env

CMD [ "/app/server" ]