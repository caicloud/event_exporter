
FROM alpine:3.7

WORKDIR /
ADD event_exporter /

ENTRYPOINT ["/event_exporter"]
CMD ["-h"]
