FROM golang:1.20-alpine AS builder

RUN apk add --no-cache make

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# RUN go build -ldflags="-s -w" -o equisplit ./src/main.go
RUN make build

FROM scratch

LABEL maintainer="Shailesh Bhimanpelli <shaileshb.0720@gmail.com> (https://shaileshhb.github.io/resume/)"
# LABEL description=""


# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/bin/equisplit", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/equisplit"]