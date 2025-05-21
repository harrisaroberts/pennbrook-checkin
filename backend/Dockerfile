# Start from the official Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (to leverage Docker cache)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the app
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]

