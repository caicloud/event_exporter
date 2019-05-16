
FROM golang:1.9-alpine as builder
WORKDIR ${GOPATH}/src/github.com/caicloud/event_exporter
COPY . ./
RUN apk add --update --no-cache make git build-base gcc abuild binutils binutils-doc gcc-doc
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/caicloud/event_exporter/event_exporter .

ENTRYPOINT ["./event_exporter"]
CMD ["-h"]
