ORG ?= Wise2c
CN ?= Wise2c CA
BITS ?= 2048
URL ?= lb-webhook.default.svc
DAYS = 36500

##help:		print this help message
.PHONY: help

help: Makefile
	@sed -n 's/^##//p' $<

.PHONY: ca.key ca.crt server.key server.csr server.crt

root-ca: ca.key ca.crt server-ca clean

server-ca: server.key server.csr server.crt

ca.key:
	@echo "generating $@"
	@openssl genrsa -out $@ ${BITS}

ca.crt:
	@echo "generating $@"
	@openssl req -x509 -new -nodes -key ca.key -subj "/CN=${CN}" -days ${DAYS} -out $@

server.key:
	@echo "generating $@"
	@openssl genrsa -out $@ ${BITS}

server.crt: server.csr
	@echo "generating $@"
	@openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
	 -CAcreateserial -out server.crt -days ${DAYS} \
	 -extensions v3_ext -extfile csr.conf

server.csr: csr.conf
	@echo "generating $@"
	@openssl req -new -key server.key -out server.csr -config csr.conf

csr.conf:
	@echo "[ req ]" >> $@
	@echo "default_bits = ${BITS}" >> $@
	@echo "prompt = no" >> $@
	@echo "default_md = sha256" >> $@
	@echo "req_extensions = req_ext" >> $@
	@echo "distinguished_name = dn" >> $@

	@echo "\n[ dn ]" >> $@
	@echo "C = CN" >> $@
	@echo "ST = ShenZhen" >> $@
	@echo "L = SZ" >> $@
	@echo "O = Wise2c" >> $@
	@echo "OU = Wise2c" >> $@
	@echo "CN = ${CN}" >> $@

	@echo "\n[ req_ext ]" >> $@
	@echo "subjectAltName = @alt_names" >> $@

	@echo "\n[ alt_names ]" >> $@
	@echo "DNS.1 = ${URL}" >> $@

	@echo "\n[ v3_ext ]" >> $@
	@echo "authorityKeyIdentifier=keyid,issuer:always" >> $@
	@echo "basicConstraints=CA:FALSE" >> $@
	@echo "keyUsage=keyEncipherment,dataEncipherment" >> $@
	@echo "extendedKeyUsage=serverAuth,clientAuth" >> $@
	@echo "subjectAltName=@alt_names" >> $@

clean:
	-@rm ca.srl server.csr csr.conf