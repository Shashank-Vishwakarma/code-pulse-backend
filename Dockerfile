# Stage 1
FROM golang:latest AS base

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /code-pulse-backend ./cmd
# /code-pulse-backend contains the compiled binary

# Stage 2
FROM alpine:latest AS run

WORKDIR /app

COPY --from=base /code-pulse-backend /app/code-pulse-backend

EXPOSE 8000

CMD ["./code-pulse-backend"]