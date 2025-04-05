# Use Go 1.23
FROM golang:1.23-alpine

WORKDIR /app
COPY . .

# Install dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o weatherapp

EXPOSE 8081

CMD ["./weatherapp"]
