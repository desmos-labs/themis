FROM golang:1.15-alpine

# Set working directory for the build
WORKDIR /go/src/github.com/desmos-labs/themis

# Add sources files
COPY . .

# Get the dependencies
RUN ["go", "mod", "tidy"]

EXPOSE 8080

# Set the entrypoint, so that the user can set the config using the CMD
ENTRYPOINT ["go", "run", "main.go"]