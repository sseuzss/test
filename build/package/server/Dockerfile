FROM debian:buster-slim

RUN mkdir -p /opt/vxserver/bin && \
    mkdir -p /opt/vxserver/data && \
    mkdir -p /opt/vxserver/logs

COPY build/package/server/preparing.sh 	/opt/vxserver/bin/
COPY build/artifacts/server/version 	/opt/vxserver/bin/
COPY build/bin/vxserver			        /opt/vxserver/bin/
COPY db/server/migrations 	            /opt/vxserver/migrations
COPY security/vconf 	                /opt/vxserver/vconf
COPY security/certs/server              /opt/vxserver/crts/server

WORKDIR /opt/vxserver

RUN chmod +x /opt/vxserver/bin/preparing.sh && \
    /opt/vxserver/bin/preparing.sh

ENTRYPOINT ["/opt/vxserver/bin/vxserver"]
