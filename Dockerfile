# BUILD
FROM golang:1.25.1-alpine AS build

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o main cmd/main.go

# PRODUCTION
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080
CMD ["./main"]
