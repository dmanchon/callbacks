# syntax=docker/dockerfile:1.2
FROM golang:1.16 as build

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o app.exe main.go

FROM scratch as bin
COPY --from=build /build/app.exe /
