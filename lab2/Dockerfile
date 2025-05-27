FROM alpine:latest

RUN apk add --no-cache bash coreutils util-linux

COPY run.sh /run.sh

RUN chmod +x /run.sh

CMD ["/bin/bash", "/run.sh"]