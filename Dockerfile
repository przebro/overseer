FROM alpine:latest
COPY ./bin/ovs-linux-amd64 /opt/overseer/ovs-linux-amd64
COPY ./config/overseer.json /etc/overseer/overseer.json
RUN ["/bin/sh"]