FROM golang:1.13-alpine as builder

ENV USER=event_exporter
ENV UID=10001
RUN adduser \
		--disabled-password \
		--gecos "" \
		--home "/nonexistent" \
		--shell "/sbin/nologin" \
		--no-create-home \
		--uid "${UID}" \
		"${USER}" && \
	apk update && apk add --no-cache git ca-certificates gcc

FROM scratch

COPY bin/event_exporter /
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER event-exporter

ENTRYPOINT ["/event_exporter"]

EXPOSE 9102
