#!/bin/bash

# Define file names for the private key and certificate
KEY_FILE="key.pem"
CERT_FILE="cert.pem"

# Generate a private key (using 2048-bit RSA)
openssl genpkey -algorithm RSA -out $KEY_FILE -pkeyopt rsa_keygen_bits:2048

# Create a configuration file for the certificate
CONFIG_FILE=$(mktemp)
cat > "$CONFIG_FILE" <<EOL
[req]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext
x509_extensions    = v3_ca
prompt             = no

[req_distinguished_name]
C  = US
ST = State
L  = City
O  = Organization
OU = OrgUnit
CN = localhost

[req_ext]
subjectAltName = @alt_names

[v3_ca]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
EOL

# Generate a self-signed certificate with SAN for localhost
openssl req -x509 -nodes -days 365 -key $KEY_FILE -out $CERT_FILE -config "$CONFIG_FILE"

# Clean up the temporary configuration file
rm "$CONFIG_FILE"

# Output the result
echo "Private key generated: $KEY_FILE"
echo "Self-signed certificate generated: $CERT_FILE"