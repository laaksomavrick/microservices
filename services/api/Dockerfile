# build stage

ARG GO_VERSION=1.11

FROM golang:${GO_VERSION}-alpine AS builder

# $GOPATH
RUN mkdir /go/src/app

# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git

# Dep
RUN go get -u github.com/golang/dep/cmd/dep

# $GOPATH
WORKDIR /go/src/app

ADD ./src ./
COPY ./Gopkg.toml  ./Gopkg.lock ./

RUN dep ensure

RUN CGO_ENABLED=0 go build -o /app .

# final stage
FROM golang:${GO_VERSION}-alpine AS final

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 3000


ENTRYPOINT ./app