# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy your Go source code into the container
COPY . .

# Build the Go application inside the container
RUN go build -o app

# Expose the port your application is running on
EXPOSE 8080

# Command to run the application
CMD ["./app"]
