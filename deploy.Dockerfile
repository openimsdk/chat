# Use an existing docker image as base
FROM ubuntu

# Set work directory
WORKDIR /chat

# Copy files from project to the work directory
COPY ./config /chat/config
COPY ./scripts /chat/scripts
COPY ./logs /chat/logs
COPY ./bin /chat/bin

# Make the script executable
RUN chmod +x ./scripts/docker_start_all.sh

# Create volumes for these directories
VOLUME ["/chat/logs"]

WORKDIR /chat/scripts


# Run the script when the container starts
CMD ["./docker_start_all.sh"]
