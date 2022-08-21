#!/bin/zsh
SERVER_CN= localhost


# Step 1 : Generate  Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout  pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -ca.crt -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private Key(server.key)
openssl genrsa -passout  pass:1111 -des3 -out server.key 4096


# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key  -out server.csr -subj "/CN=${SERVER_CN}"


# Step 4 : Signi  the certificate  with the CA  created
openssl x509 -req -passin pass:1111 -days 365  -in server.csr  -CA ca.crt -CAkey ca.key  -set_serial 01 -out server.crt