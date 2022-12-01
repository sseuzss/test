FROM debian:buster-slim

RUN mkdir -p /opt/api/logs && \
    mkdir -p /opt/api/logsssl

COPY build/bin/vxapi                   /opt/api/bin/
COPY db/api/migrations                 /opt/api/migrations
COPY build/package/api/templates       /opt/api/templates
COPY security/certs/api                /opt/api/crts
COPY build/package/api/entrypoint.sh   /opt/api/bin/

RUN chmod +x /opt/api/bin/entrypoint.sh

WORKDIR /opt/api

ENTRYPOINT ["/opt/api/bin/entrypoint.sh", "/opt/api/bin/vxapi"]
