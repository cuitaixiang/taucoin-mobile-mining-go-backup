# Build Gtau in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-tau
RUN cd /go-tau && make gtau

# Pull Gtau into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-tau/build/bin/gtau /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["gtau"]