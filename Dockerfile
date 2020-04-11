FROM gcr.io/distroless/static

COPY event_exporter /

USER nobody

ENTRYPOINT ["/event_exporter"]

EXPOSE 9102
