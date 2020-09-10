FROM alpine:3.12

RUN apk --no-cache add \
    bash \
    zip \
    curl \
    git

ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "help" ]

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY cider /bin/cider
