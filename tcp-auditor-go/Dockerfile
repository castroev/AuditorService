# build go
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git gcc g++ libc-dev musl-dev ca-certificates bash
RUN addgroup -S appgroup && adduser -S -D -H -h /app -G appgroup appuser
# COPY . $GOPATH/src/bitbucket.tylertech.com/scm/spy/tcp-auditor
COPY . $GOPATH/src/bitbucket.tylertech.com/scm/spy/tcp-auditor
WORKDIR $GOPATH/src/bitbucket.tylertech.com/scm/spy/tcp-auditor

ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags='-linkmode external -extldflags "-static"' -a -o /go/bin/tcp-auditor .

# build final image
FROM scratch
EXPOSE 8080
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --chown=appuser:appgroup --from=builder /go/bin/tcp-auditor /app/
USER appuser
WORKDIR /app
ENV GIN_MODE=release
ENTRYPOINT ["/app/tcp-auditor"]
