FROM golang:1.24.5-alpine As builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bot .

CMD [ "./bot" ]