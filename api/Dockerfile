# Official Golang image https://hub.docker.com/_/golang
FROM golang:1.23

# Working directory
WORKDIR /app

# Copy everything at /app
COPY . /app

# Build the go app
RUN go build -o main .

# Define the command to run the app
CMD ["./main"]
