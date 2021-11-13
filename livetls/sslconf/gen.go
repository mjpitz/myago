package main

// Generate ca_key.pem and ca.pem components
//go:generate openssl genrsa -out ca_key.pem 4096
//go:generate openssl req -new -x509 -key ca_key.pem -days 1 -out ca.pem -config ca.conf

// Generate key.pem and cert.csr
//go:generate openssl genrsa -out key.pem 4096
//go:generate openssl req -new -key key.pem -out cert.csr -config cert.conf

// Sign cert.csr using the ca and output to cert.pem
//go:generate openssl x509 -req -in cert.csr -CA ca.pem -CAkey ca_key.pem -CAcreateserial -days 1 -out cert.pem -extfile cert.conf
