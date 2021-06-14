FROM docker.io/golang:1.16-alpine as build

# Create appuser and group.
ARG USER=bw
ARG GROUP=bw
ARG USER_ID=10001
ARG GROUP_ID=10001

ENV CGO_ENABLED=0

# Prepare build stage - can be cached
WORKDIR /work
RUN apk add -U --no-cache \
        make protoc gcc musl-dev && \
    addgroup -S -g "${GROUP_ID}" "${GROUP}" && \
    adduser \
        --disabled-password \
        --gecos "" \
        --home "/nonexistent" \
        --shell "/sbin/nologin" \
        --no-create-home \
        -G "${GROUP}" \
        --uid "${USER_ID}" \
        "${USER}"

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download && \
    go get -u github.com/golang/mock/mockgen@latest

COPY ./ ./

RUN GOOS=linux GOARCH=amd64 go build -o bw-crowdedness -ldflags='-w -s' -a -installsuffix cgo ./

FROM docker.io/alpine:3.13

WORKDIR /app

COPY --from=build /etc/passwd /etc/group /etc/
COPY --from=build --chown=$USER_ID:$GROUP_ID /work/bw-crowdedness ./

EXPOSE 9091

USER $USER

ENTRYPOINT ["/app/bw-crowdedness"]