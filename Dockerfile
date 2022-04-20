FROM golang:bullseye
WORKDIR /src
COPY go.mod go.sum main.go ./
RUN go build

FROM debian:bullseye-slim
WORKDIR /app
RUN apt-get update && apt-get install -y tini && rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["/usr/bin/tini", "--"]
COPY --from=0 /src/quic-go-server-example /app/
CMD ["/app/quic-go-server-example"]
