build-docker: build-for-linux
	docker build -t hup-on-notify-example .

build-for-linux: *.go
	GOOS=linux go build -o hup-on-notify-linux
