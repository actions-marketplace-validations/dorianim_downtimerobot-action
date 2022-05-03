FROM golang:1.16.0-buster AS builder
LABEL stage=intermediate
COPY . /downtimerobot
WORKDIR /downtimerobot
ENV GO111MODULE=on
RUN apt update && apt install -y ca-certificates && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 && \
    go test ./... && \
    go build -a --installsuffix cgo -v -tags netgo -ldflags '-extldflags "-static"' -o /downtimerobot .

FROM scratch
LABEL maintainer="Dorian Zedler <dev@dorian.im>"
WORKDIR /
COPY --from=builder \
    /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /downtimerobot ./
ENTRYPOINT [ "/downtimerobot" ]