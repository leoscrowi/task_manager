FROM golang:1.24.2

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o task-service-server ./cmd

EXPOSE 8080

CMD ["./task-service-server"]