#!/bin/bash
#generate keys for rootca server and client
openssl genrsa -out ./docker/cert/root_ca.key 4096
openssl genrsa -out ./docker/cert/mongodb.key 2048
openssl genrsa -out ./docker/cert/overseer_srv.key 2048
openssl genrsa -out ./docker/cert/client.key 2048

#create root certficiate
openssl req -x509 -new -nodes -key ./docker/cert/root_ca.key -sha256 -days 1024 \
-out ./docker/cert/root_ca.crt \
-subj "/C=US/ST=State/L=City/O=overseer/OU=DEV/CN=localhost/emailAddress=localuser@localhost.com"
#create signing request
openssl req -new -key ./docker/cert/mongodb.key -out ./docker/cert/mongodb.csr \
-subj "/C=US/ST=State/L=City/O=ovsmongodb/OU=DEV/CN=localhost/emailAddress=localuser@localhost.com"
#create certificate for mongodb server
openssl x509 -req -in ./docker/cert/mongodb.csr -CA ./docker/cert/root_ca.crt -CAkey ./docker/cert/root_ca.key \
-CAcreateserial -out ./docker/cert/mongodb.crt -days 1024 -sha256

#create signing request for overseer scheduler
openssl req -new -key ./docker/cert/overseer_srv.key -out ./docker/cert/overseer_srv.csr \
-subj "/C=US/ST=State/L=City/O=overseer_srv/OU=DEV/CN=localhost/emailAddress=localuser@localhost.com" \
-addext "subjectAltName=DNS:*.domain.com,IP:127.0.0.1"
#create certificate for overseer scheduler
openssl x509 -req -in ./docker/cert/overseer_srv.csr -CA ./docker/cert/root_ca.crt -CAkey ./docker/cert/root_ca.key \
-CAcreateserial -out ./docker/cert/overseer_srv.crt -days 1024 -sha256

#create certificate for client
openssl req -new -key ./docker/cert/client.key -out ./docker/cert/client.req \
-subj "/C=US/ST=State/L=City/O=overseer/OU=DEV/CN=localhost/emailAddress=localuser@localhost.com"
openssl x509 -req -in ./docker/cert/client.req -CA ./docker/cert/root_ca.crt -CAkey ./docker/cert/root_ca.key \
-set_serial 101010 -extensions client -days 1024 -out ./docker/cert/client.crt

cat ./docker/cert/mongodb.key ./docker/cert/mongodb.crt > ./docker/cert/mongodb.pem

#Subject: C = PL, ST = State, L = Lodz, O = Couchdrv, OU = section, CN = localuser, emailAddress = localuser@localhost.com
#Subject: C = US, ST = State, L = City, O = CouchDB driver ltd., OU = DEV, CN = localuser, emailAddress = localuser@localhost.com

#Subject: emailAddress = user@host.com, C = US, ST = State, L = City, O = CouchDB driver ltd., OU = DEV, CN = localhost
#Subject: C = PL, ST = State, L = City, O = Company ltd, OU = section, CN = localuser, emailAddress = user@localhost.com