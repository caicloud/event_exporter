FROM debian:stretch-curl

COPY bin/event_exporter /

USER nobody

ENTRYPOINT ["/event_exporter"]

EXPOSE 9102
