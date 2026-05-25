FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/service ./backend/cmd/investment-service/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/service /app/service
ENTRYPOINT ["/app/service"]
