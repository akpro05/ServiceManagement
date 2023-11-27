FROM debian:stable-slim

# Install CA certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY ServiceManagement /
COPY config.yml /
RUN mkdir /templates
COPY templates /templates

RUN chmod +x /ServiceManagement

ENTRYPOINT ["/ServiceManagement"]
