FROM golang:1.15-alpine as builder
ENV GOFLAGS -mod=vendor

WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY ./cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./...

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /go/bin/cmd /usr/local/bin/leader-election
ENTRYPOINT ["leader-election"]
