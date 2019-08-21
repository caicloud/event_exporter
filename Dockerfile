
FROM golang:1.12-alpine as builder
WORKDIR ${GOPATH}/src/github.com/caicloud/event_exporter
RUN apk add --update --no-cache make git build-base gcc abuild binutils binutils-doc gcc-doc
COPY . ./
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/caicloud/event_exporter/event_exporter .

ENTRYPOINT ["./event_exporter"]
