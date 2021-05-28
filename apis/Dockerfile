FROM golang:1.15-alpine

# Set working directory for the build
WORKDIR /go/src/github.com/desmos-labs/themis/apis

# Add sources discord
COPY . .

# Get the dependencies
RUN ["go", "mod", "tidy"]

# Set the entrypoint, so that the user can set the config using the CMD
ENTRYPOINT ["go", "run", "main.go"]