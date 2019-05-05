ARG GO_VERSION=1.11

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git

WORKDIR /
COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/grpc-server cmd/server/account/main.go

FROM alpine:latest
COPY --from=builder /bin/grpc-server /bin/grpc-server

EXPOSE 8080
ENTRYPOINT ["./bin/grpc-server"]
