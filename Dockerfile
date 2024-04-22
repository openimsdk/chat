# Use Go 1.21 as the base image for building the application
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /openim-chat

# Set the Go proxy to improve dependency resolution speed
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download
# Copy all files from the current directory into the container
COPY . .


RUN go install github.com/magefile/mage@latest

RUN mage build

# Use Alpine Linux as the final base image due to its small size and included utilities
FROM alpine:latest

# Install necessary packages, such as bash, to ensure compatibility and functionality
RUN apk add --no-cache bash

# Copy the compiled binaries and mage from the builder image to the final image
COPY --from=builder /openim-chat/_output /openim-chat/_output
COPY --from=builder /go/bin/mage /usr/local/bin/mage

# Set the working directory to /openim-chat within the container
WORKDIR /openim-chat

# Set up volume mounts for the configuration directory and logs directory
VOLUME ["/openim-chat/config", "/openim-chat/_output/logs"]

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]


