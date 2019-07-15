FROM golang:1.11-alpine as builder
ARG gitCommit

COPY ./ /go/src/github.com/mojo-zd/kube-webhooks
WORKDIR /go/src/github.com/mojo-zd/kube-webhooks/cmd/lb-webhook

RUN CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
    -ldflags "-s -w -X main.gitCommit=${gitCommit}" -o ./lb-webhook && \
    mkdir -p /app/go/bin/ && \
    mv ./lb-webhook /app/go/bin && \
    chmod +x /app/go/bin/lb-webhook

FROM alpine:3.9.2
COPY --from=builder /app/ /
EXPOSE 443
ENTRYPOINT ["/go/bin/lb-webhook"]