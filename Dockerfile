# For local development
FROM golang:1.21-bullseye

RUN apt-get update && \
    apt-get install -y gcc g++ bash && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install delve and air for debugging and hot reload
# WARN: The built image will be bigger than 1GB!
RUN go install github.com/go-delve/delve/cmd/dlv@v1.22.0 && \
    go install github.com/cosmtrek/air@v1.44.0

WORKDIR /app
COPY . .
RUN go mod download

CMD ["air", "-c", ".air.toml"]
