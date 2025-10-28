# BUILD
FROM golang:1.25.1-alpine AS build

RUN apk add --no-cache git curl

# Install typst using the official Typst release from GitHub
RUN mkdir -p /tmp/typst-install && \
    curl -fsSL https://github.com/typst/typst/releases/download/v0.14.0/typst-x86_64-unknown-linux-musl.tar.xz \
    -o /tmp/typst-install/typst.tar.xz && \
    tar -xf /tmp/typst-install/typst.tar.xz -C /tmp/typst-install

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o main cmd/main.go

# PRODUCTION
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=build /tmp/typst-install/typst-x86_64-unknown-linux-musl/typst /usr/local/bin/typst
RUN chmod +x /usr/local/bin/typst

WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/main .

# Copy Typst templates
COPY templates templates/

EXPOSE 8080
CMD ["./main"]
