FROM debian:buster-slim

RUN mkdir -p /opt/api/{logs,ssl}

COPY ../../bin/vxapi                   /opt/api/bin/
COPY ../../../db/api/migrations               /opt/api/migrations
COPY templates     /opt/api/templates
COPY ../../../security/certs/api         /opt/api/crts
COPY entrypoint.sh /opt/api/bin/

RUN chmod +x /opt/api/bin/entrypoint.sh

WORKDIR /opt/api

ENTRYPOINT ["/opt/api/bin/entrypoint.sh", "/opt/api/bin/vxapi"]
