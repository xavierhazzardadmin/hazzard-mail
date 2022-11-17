FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o /tmp/blog ./main.go

FROM alpine:edge

COPY --from=build /tmp/blog /sbin/blog

CMD /sbin/blog
