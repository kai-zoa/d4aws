FROM alpine:3.5

LABEL io.whalebrew.config.volumes '["~/.aws:/home/d4aws/.aws:ro"]'

ARG user=d4aws
ARG group=d4aws
ARG uid=1000
ARG gid=1000

RUN apk add --no-cache ca-certificates \
 && rm -rf /var/cache/apk/* \
 && addgroup -g ${gid} ${group} \
 && adduser -h /home/d4aws -u ${uid} -G ${group} -s /bin/sh -D ${user}

VOLUME /home/d4aws/.aws

COPY d4aws /bin/d4aws

USER ${user}

ENTRYPOINT ["d4aws"]
