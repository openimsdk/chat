# Use Ubuntu 18.04 with Go installed as the base image for building the application
FROM ubuntu:18.04 as builder

# Install Go and necessary tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    golang-go \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Define the base directory for the application as an environment variable
ENV SERVER_DIR=/openim-chat \
    GOPATH=/go \
    PATH=$GOPATH/bin:$PATH

# Set the working directory inside the container based on the environment variable
WORKDIR $SERVER_DIR

# Set the Go proxy to improve dependency resolution speed
ENV GOPROXY=https://goproxy.io,direct

# Copy all files from the current directory into the container
COPY . .

# Download dependencies
RUN go mod download

# Install Mage to use for building the application
RUN go get -u github.com/magefile/mage

# Optionally build your application if needed
# RUN mage build

# Using Ubuntu 18.04 for the final image
FROM ubuntu:18.04

# Install necessary packages, such as bash
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Set the environment and work directory
ENV SERVER_DIR=/openim-chat
WORKDIR $SERVER_DIR

# Copy the compiled binaries and mage from the builder image to the final image
COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/
COPY --from=builder $SERVER_DIR/config $SERVER_DIR/
COPY --from=builder /go/bin/mage /usr/local/bin/mage
COPY --from=builder $SERVER_DIR/magefile_windows.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile_unix.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/start-config.yml $SERVER_DIR/
COPY --from=builder $SERVER_DIR/go.mod $SERVER_DIR/
COPY --from=builder $SERVER_DIR/go.sum $SERVER_DIR/

# Set up volume mounts for the configuration directory and logs directory
VOLUME ["$SERVER_DIR/config", "$SERVER_DIR/_output/logs"]

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]
