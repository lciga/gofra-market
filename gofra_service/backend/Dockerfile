### Builder stage: use official Go image to compile the app with module cache
FROM golang:1.25-alpine AS builder

# small utilities for build & git-based module fetches
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# copy only go.mod/go.sum to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy the rest and build
COPY . .

# Build statically (disable cgo), optimized flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /gofra ./cmd/api

### Final image: small runtime image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /gofra /app/gofra

# non-root user
RUN addgroup -S app && adduser -S -G app app
USER app

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["/app/gofra"]
