# Use official Golang image
FROM golang:1.20 

# Set environment variables
WORKDIR /app

# Copy Go module files first for better caching
COPY go.mod ./

# Download dependencies (prevents re-downloading on code changes)
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Install required packages
RUN go get -u github.com/gorilla/mux

# Build the application
RUN go build -o main main.go

# Expose port 8080
EXPOSE 8080

# Start the application
CMD ["./main"]