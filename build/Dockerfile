FROM debian:stretch-slim

COPY bin/event_exporter /

USER nobody

ENTRYPOINT ["/event_exporter"]

EXPOSE 9102
