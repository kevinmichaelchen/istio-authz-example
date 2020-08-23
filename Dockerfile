#syntax=docker/dockerfile-upstream:master-experimental
FROM golang:1.13 AS builder

ENV GO111MODULE on
ENV GOPRIVATE github.com/kevinmichaelchen

ARG GITHUB_USER
ARG GITHUB_ACCESS_TOKEN

COPY clone-deps.sh go.mod go.sum /go/app/
WORKDIR /go/app

# Install ssh client and git
RUN apt-get update && apt-get install -y openssh-client git

# Download dependencies
RUN --mount=type=ssh ./clone-deps.sh ${GITHUB_USER} ${GITHUB_ACCESS_TOKEN}

COPY . /go/app
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest as app

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 \
 && wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe

COPY --from=builder /go/app/app /app/app

RUN apk add --no-cache ca-certificates

WORKDIR /app
CMD ["./app"]