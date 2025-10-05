FROM golang:alpine AS builder

ARG VERSION=dev

WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -ldflags="-s -w \
     -X 'github.com/ipsums9527/monitor-etl/config.version=$VERSION' \
     -X 'github.com/ipsums9527/monitor-etl/config.date=$(date -Iseconds)' \
     " -trimpath -v -o monitor-etl main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/monitor-etl .
ENTRYPOINT ["/app/monitor-etl"]
