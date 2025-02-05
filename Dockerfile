FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk update &&\
    apk add --no-cache make git upx ca-certificates &&\
    update-ca-certificates &&\
    adduser -SHu 10001 scratch
ENV CGO_ENABLED=0
COPY go.mod .
RUN go mod download
COPY . .
RUN \
  go build \
    -trimpath \
    -ldflags='-s -w' \
    -o go-vanity \
    . \
  && \
  upx -5 go-vanity

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/go-vanity /go-vanity
USER scratch
CMD ["/go-vanity"]
EXPOSE 5000
