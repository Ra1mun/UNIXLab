FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o main .

CMD ["./main"]