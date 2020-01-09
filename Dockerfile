ARG GO_VERSION=1.13.5
ARG ALPINE_VERSION=3.11
ARG BASE_IMAGE=scratch

############################
# STEP 1 build executable binary
############################
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN echo 'hosts: files dns' > /etc/nsswitch.conf.build

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
# tzdata is time zone stuff
# bash, make, curl are optional
RUN apk add --no-cache --update \
        tzdata \
        openssh-client \
        git \
        ca-certificates \
        bash make curl \
        && update-ca-certificates --fresh \
        && adduser -D -g '' appuser
ENV GOFLAGS="-mod=readonly"

# Moving outside of $GOPATH forces modules on without having to set ENVs
WORKDIR /src

# Add go.mod and go.sum first to maximize caching
COPY ./go.mod ./go.sum ./

# Copy in the project
COPY . .

# Compile...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/main ./main.go
RUN ls -la /go/bin/
############################
# STEP 2 build a small image
############################
FROM ${BASE_IMAGE}

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/main /go/bin/main

# Use an unprivileged user.
USER appuser

# Run the hello binary.
ENTRYPOINT ["/go/bin/main"]