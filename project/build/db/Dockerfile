FROM mysql:8.0

COPY build/db/conf.d/ /etc/mysql/conf.d/
COPY build/db/initdb.d/ /docker-entrypoint-initdb.d/
