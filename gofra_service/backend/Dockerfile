FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /gofra ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /gofra /app/gofra

RUN addgroup -S app && adduser -S -G app app
USER app

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["/app/gofra"]
