help:
	echo "commands: login, build, push, release;"
	echo "login user=aaa@bbb.com pwd=123 [server="https://index.docker.io/v1"]; "
	echo "build [tag=0.1]; "
	echo "push [tag=0.1]\n"
login:
	docker login -u="$(user)" -p="$(pwd)" $(or "$(server)", "https://index.docker.io/v1")
build:
	mkdir -p dist
	go get
	GOOS=linux GOARCH=amd64  go build -a -o dist/bitcoin -work cmd/bitcoin/main.go
	docker build -t s1rxploit/bitcoinwallet .
	#docker tag ${ID} s1rxploit/bitcoinwallet:$(or $(tag), "latest")
	#docker tag ${ID} s1rxploit/bitcoinwallet:latest
push:
	docker push OpenBazaar/bitcoinwallet:latest
release: build push
