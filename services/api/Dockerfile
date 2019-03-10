# build stage

ARG GO_VERSION=1.11

FROM golang:${GO_VERSION}-alpine AS builder

# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

ENV GO111MODULE=on

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -o /app .

# final stage
FROM golang:${GO_VERSION}-alpine AS final

WORKDIR /app

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /app /app

EXPOSE 3000

USER nobody:nobody

ENTRYPOINT ./app