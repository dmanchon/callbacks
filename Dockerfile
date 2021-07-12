# syntax=docker/dockerfile:1
FROM golang:1.16

WORKDIR /build
COPY go.mod go.sum .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app.exe main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /build/app.exe .
