FROM index.caicloud.io/debian:jessie
MAINTAINER zhoushaolei <shaolei@caicloud.io>

WORKDIR /
ADD event_exporter /

# Set the timezone to Shanghai
RUN echo "Asia/Shanghai" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata && \
    sed -i "s/httpredir.debian.org/mirrors.163.com/g" /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates

ENTRYPOINT ["/event_exporter"]
CMD ["-h"]
