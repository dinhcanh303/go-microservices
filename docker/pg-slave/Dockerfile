FROM postgres:14-alpine
RUN chmod 777 /tmp
RUN apk update && apk add iputils-ping openssh
COPY ./docker/pg-slave/setup-slave.sh /setup-slave.sh
RUN chmod +x /setup-slave.sh
ENTRYPOINT ["/setup-slave.sh"]
CMD ["gosu","postgres","postgres"]