############################
# Run from parent directory
# docker build -f api/Dockerfile .
############################

############################
# STEP 1 build executable binary
############################
FROM golang:1.12-alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk add --no-cache git make gcc musl-dev linux-headers

WORKDIR $GOPATH/src/github.com/bahelit/confirmerator/

# COPY . .
# Fetch dependencies.
# Using go get.
# RUN go get -d -v
# Using go mod.
# RUN go mod download
# Using dep.
# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY ../../cmd/confirmerator-api ./

# Build the binary
RUN cd api
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/confirmerator-api
############################
# STEP 2 build a small image
############################
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Copy our executable
COPY --from=builder /go/bin/confirmerator-api /go/bin/confirmerator-api

# Open the pod bay doors.
EXPOSE 80

# Run the binary.
ENTRYPOINT ["/go/bin/confirmerator-api"]