#!/usr/bin/env bash

set -e

# gencerts.sh generates the certificates for the webhook authz plugin tests.
#
# It is not expected to be run often (there is no go generate rule), and mainly
# exists for documentation purposes.

cat > server.conf << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
IP.1 = 127.0.0.1
EOF

cat > client.conf << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
EOF

# Create a certificate authority
openssl genrsa -out caKey.pem 2048
openssl req -x509 -new -nodes -key caKey.pem -days 100000 -out caCert.pem -subj "/CN=webhook_imagepolicy_ca"

# Create a second certificate authority
openssl genrsa -out badCAKey.pem 2048
openssl req -x509 -new -nodes -key badCAKey.pem -days 100000 -out badCACert.pem -subj "/CN=webhook_imagepolicy_ca"

# Create a server certiticate
openssl genrsa -out serverKey.pem 2048
openssl req -new -key serverKey.pem -out server.csr -subj "/CN=webhook_imagepolicy_server" -config server.conf
openssl x509 -req -in server.csr -CA caCert.pem -CAkey caKey.pem -CAcreateserial -out serverCert.pem -days 100000 -extensions v3_req -extfile server.conf

# Create a client certiticate
openssl genrsa -out clientKey.pem 2048
openssl req -new -key clientKey.pem -out client.csr -subj "/CN=webhook_imagepolicy_client" -config client.conf
openssl x509 -req -in client.csr -CA caCert.pem -CAkey caKey.pem -CAcreateserial -out clientCert.pem -days 100000 -extensions v3_req -extfile client.conf

outfile=certs_test.go

cat > $outfile << EOF
// This file was generated using openssl by the gencerts.sh script
// and holds raw certificates for the imagepolicy webhook tests.

//lint:file-ignore U1000 Ignore all unused code, it's generated

package imagepolicy
EOF

for file in caKey caCert badCAKey badCACert serverKey serverCert clientKey clientCert; do
	data=$(cat ${file}.pem)
	echo "" >> $outfile
	echo "var $file = []byte(\`$data\`)" >> $outfile
done

# Clean up after we're done.
rm ./*.pem
rm ./*.csr
rm ./*.srl
rm ./*.conf
