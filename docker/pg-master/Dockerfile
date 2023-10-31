FROM postgres:14-alpine
RUN chmod 777 /tmp
RUN apk update && apk add iputils-ping openssh
COPY ./docker/pg-master/setup-master.sh /docker-entrypoint-initdb.d/setup-master.sh
RUN chmod 777 /docker-entrypoint-initdb.d/setup-master.sh

COPY ./docker/pg-master/setup-entry.sh /setup-entry.sh
RUN chmod 777 /setup-entry.sh
CMD ["./setup-entry.sh"]