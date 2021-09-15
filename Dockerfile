FROM golang:1.16.6-alpine AS builder
RUN mkdir /build
ADD go.mod go.sum main.go /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/hazzard-mail /app/
WORKDIR /app
CMD ["./hazzard-mail"]
