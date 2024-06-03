# Example:
# docker build .
# docker run --name sonic --entrypoint sonictool 43eedf15b4d0 --datadir=/var/sonic genesis fake 1
# docker run --volumes-from sonic -p 5050:5050 -p 5050:5050/udp -p 18545:18545 43eedf15b4d0 --http --http.addr=0.0.0.0

FROM golang:1.22-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/Sonic
COPY . .

ARG GOPROXY
RUN make all


FROM alpine:latest

COPY --from=builder /go/Sonic/build/sonicd /usr/local/bin/
COPY --from=builder /go/Sonic/build/sonictool /usr/local/bin/

EXPOSE 18545 18546 5050 5050/udp

VOLUME /var/sonic

ENTRYPOINT ["sonicd", "--datadir=/var/sonic"]
