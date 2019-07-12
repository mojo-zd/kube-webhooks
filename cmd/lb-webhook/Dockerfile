FROM golang:1.11-alpine as builder
ARG gitCommit

COPY ./ /go/src/gitee.com/wisecloud/kube-webhooks
WORKDIR /go/src/gitee.com/wisecloud/kube-webhooks/cmd/lb-webhook

RUN echo "https://mirrors.aliyun.com/alpine/v3.9/main/"      >  /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.9/community/" >> /etc/apk/repositories && \
    apk update && \
    apk add -U tzdata && \
    CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
    -ldflags "-s -w -X main.gitCommit=${gitCommit}" -o ./lb-webhook && \
    mkdir -p /app/go/bin/ && \
    mv ./lb-webhook /app/go/bin && \
    chmod +x /app/go/bin/lb-webhook

FROM alpine:3.9.2
COPY --from=builder /app/ /
EXPOSE 443
ENTRYPOINT ["/go/bin/lb-webhook"]