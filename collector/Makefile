.PHONY: all
all: start

.PHONY: certs
certs: server.key server.crt
server.key server.crt:
	openssl req -x509 -newkey rsa:4096 \
		-keyout server.key \
		-out server.crt -days 3650 \
		-subj '/CN=localhost' \
		-nodes

.PHONY: start
start: 
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down -v --remove-orphans