build-docker: build-for-linux
	docker build -t hup-on-notify-example .

build-for-linux: *.go
	go fmt
	GOOS=linux go build -o hup-on-notify-linux

clean:
	rm hup-on-notify*
