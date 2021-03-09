FROM alpine:3.12.0

RUN apk --no-cache add ca-certificates=20191127-r4

WORKDIR /service

ARG RUN_USER=service

RUN adduser -S -D -H -u 1001 -s /sbin/nologin -G root -g $RUN_USER $RUN_USER

COPY service .

RUN chgrp -R 0 /service && chmod -R g+rX /service

USER $RUN_USER

CMD ./service