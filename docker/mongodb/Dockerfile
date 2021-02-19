FROM mongo
WORKDIR $PWD
ADD ./docker/cert/*.crt /opt/mongo/cert/
ADD ./docker/cert/*.key /opt/mongo/cert/
ADD ./docker/cert/*.csr /opt/mongo/cert/
ADD ./docker/cert/*.req /opt/mongo/cert/
ADD ./docker/cert/*.pem /opt/mongo/cert/
COPY ./docker/mongodb/first-run.js /docker-entrypoint-initdb.d/
