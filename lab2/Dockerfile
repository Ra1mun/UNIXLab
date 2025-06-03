FROM alpine:latest

RUN apk add --no-cache bash dos2unix

COPY run.sh /run.sh

RUN mkdir -p /shared && \
    dos2unix /run.sh && \
    chmod +x /run.sh

CMD ["/bin/bash", "/run.sh"]