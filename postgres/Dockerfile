FROM postgres:alpine

COPY create-db.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/create-db.sh

EXPOSE 5432
