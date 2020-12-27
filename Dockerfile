# hadolint ignore=DL3007
FROM alpine:latest

COPY solution-sample.sh /usr/bin/solution-sample.sh

RUN chmod +x /usr/bin/solution-sample.sh

ENTRYPOINT [ "/usr/bin/solution-sample.sh" ]
