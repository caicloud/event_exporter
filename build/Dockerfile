FROM debian:stretch-slim

COPY bin/event_exporter /

USER 65534

ENTRYPOINT ["/event_exporter"]

EXPOSE 9102
