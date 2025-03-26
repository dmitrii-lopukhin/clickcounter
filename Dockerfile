FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o clicks-counter ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/clicks-counter .

EXPOSE 8081

CMD ["/app/clicks-counter"]
