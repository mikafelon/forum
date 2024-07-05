# Use the official Golang image as a parent image
FROM golang:1.19

# Set the working directory in the container
WORKDIR /usr/src/app

# Add metadata to the image
LABEL version="1.0"
LABEL description="Forum Application"
LABEL author="Your Name"

# Set CGO_ENABLED to 1
ENV CGO_ENABLED=1
ENV PORT=8080

# Copy the go.mod and go.sum files to the working directory
COPY go.mod ./

# Download and install dependencies
RUN go mod download && go mod tidy

# Copy the rest of the application code to the working directory
COPY . .

# Build the application with CGO enabled
RUN GOOS=linux go build -a -installsuffix cgo -o app .

# Expose port 443 for the application
EXPOSE 8080

# Run the application
CMD ["./app"]