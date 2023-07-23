# Use the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy all the Go source files to the container
COPY *.go ./

# Install the dependencies needed by your application
RUN go get -d -v ./...

# Build the Go application inside the container
RUN go build -o mighty-saver-rabbit

# Expose the port that your application listens on
EXPOSE 8080

# Define the command to run your application
CMD [ "./mighty-saver-rabbit" ]
