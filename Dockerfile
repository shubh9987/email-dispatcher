FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o email-dispatcher .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/email-dispatcher .
COPY --from=builder /app/email.tmpl .
COPY --from=builder /app/emails.csv .

EXPOSE 8080

CMD ["./email-dispatcher"]