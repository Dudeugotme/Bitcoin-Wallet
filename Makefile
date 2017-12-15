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
	GOOS=linux GOARCH=amd64 go build -a -o dist/bitcoin -work cmd/bitcoin/main.go
	docker build -t opendomido/bitcoinwallet .
	#docker tag ${ID} opendomido/bitcoinwallet:$(or $(tag), "latest")
	#docker tag ${ID} opendomido/bitcoinwallet:latest
push:
	docker push opendomido/bitcoinwallet:latest
release: build push
	