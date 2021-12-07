# Builder

ARG GITHUB_PATH=github.com/ozonmp/cnm-purchase-api

FROM golang:1.16-alpine AS builder

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git protoc protobuf protobuf-dev curl
COPY Makefile Makefile
RUN make deps-go
COPY . .
RUN make build-go

# gRPC Server

FROM alpine:latest as server
rUN apk add --update make git protoc protobuf protobuf-dev curl
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/grpc-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations

RUN chown root:root grpc-server

EXPOSE 50051
EXPOSE 8080
EXPOSE 9101

CMD ["./grpc-server"]
