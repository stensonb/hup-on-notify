FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add squid && \
    chown squid:squid /run && \
    addgroup squid tty && \
    chown squid:squid /etc/squid/squid.conf

COPY hup-on-notify-linux /usr/local/bin/hup-on-notify

USER squid

CMD ["squid", "-N", "-d", "1"]
