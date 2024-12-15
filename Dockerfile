FROM golang:1.23.4

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify
COPY . .

RUN go build -o app ./cmd/main.go

CMD ["./app"]
