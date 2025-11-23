FROM golang:1.24.4-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr_reviewer_app ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata bash

WORKDIR /app

COPY --from=builder /app/pr_reviewer_app .

ENV DATABASE_URL=postgres://user:password@db:5432/pr_reviewer?sslmode=disable

CMD ["./pr_reviewer_app"]
