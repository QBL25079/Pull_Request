
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git bash postgresql-client
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pr_reviewer_app ./cmd/server/main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client bash

WORKDIR /app

COPY --from=builder /pr_reviewer_app /pr_reviewer_app
COPY --from=builder /app/wait-for-it.sh /wait-for-it.sh

COPY migrations migrations

RUN chmod +x /wait-for-it.sh /pr_reviewer_app
