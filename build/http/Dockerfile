FROM golang:1.20.3 AS builder

WORKDIR /builder
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Build Go App
RUN CGO_ENABLED=0 GOOS=linux go build -o e-resep-service .

FROM alpine:3.11.3

WORKDIR /app

COPY --from=builder ./builder/e-resep-service .
COPY --from=builder ./builder/.env .

# Expose port
EXPOSE 8888

# Command to run the executeable
ENTRYPOINT ["./e-resep-service","http"]
