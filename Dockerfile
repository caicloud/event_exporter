FROM index.caicloud.io/debian:jessie
MAINTAINER zhoushaolei <shaolei@caicloud.io>

WORKDIR /
ADD event_exporter /

# Set the timezone to Shanghai
RUN echo "Asia/Shanghai" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata

ENTRYPOINT ["/event_exporter"]
CMD ["-h"]
