.PHONY: build generate_certs clean

CA_KEY := "ca.key"
CA_CSR := "ca.csr"
SERVER_CERT := "server.crt"
EXT_FILE := "self-signed-cert.ext"
BINARY := "tpm"

all := $(CA_KEY) $(CA_CSR)$(SERVER_CERT) build

build: $(SERVER_CERT)
	go build -ldflags="-s -w" -o $(BINARY) main.go

$(SERVER_CERT): generate_certs
	@echo "subjectAltName = @alt_names" > $(EXT_FILE)
	@echo "[alt_names]" 								>> $(EXT_FILE)
	@echo "DNS.1 = localhost" 					>> $(EXT_FILE)
	@openssl x509 -req -days 365 \
		-in $(CA_CSR) \
		-signkey $(CA_KEY) \
		-out $(SERVER_CERT) \
		-extfile $(EXT_FILE)

generate_certs:
	@openssl req \
		-newkey rsa:2048 -nodes -keyout $(CA_KEY) \
		-new -subj "/C=KR/CN=localhost" \
		-out $(CA_CSR)

clean:
	@rm -rf $(CA_KEY) $(CA_CSR) $(SERVER_CERT) $(EXT_FILE) $(BINARY)

