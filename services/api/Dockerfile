# -- build stage
ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}-alpine AS builder
# $GOPATH
RUN mkdir /go/src/app
# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git
# Dep
RUN go get -u github.com/golang/dep/cmd/dep
# set workdir to $GOPATH base
WORKDIR /go/src/app
# Copy dependencies map and src to container
COPY ./Gopkg.toml ./Gopkg.lock ./src ./
# Install dependencies
RUN dep ensure
# Build the app
RUN CGO_ENABLED=0 go build -o /app .

# -- app stage
FROM alpine AS final
# Set workdir to /app in container
WORKDIR /app
# Copy binary
COPY --from=builder /app /app
# Expose port
EXPOSE 3000
# Run app
ENTRYPOINT ./app