FROM golang:1.22.3

WORKDIR /app

RUN apt-get update && apt-get install -y netcat-openbsd
COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main .

WORKDIR /app

EXPOSE 8080

CMD ["./main"]
