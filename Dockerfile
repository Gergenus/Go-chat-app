FROM golang:1.22-alpine as builder

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o ./chat-app ./cmd

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/chat-app ./chat-app
COPY --from=builder /build/.env .env
CMD ["/app/chat-app"]