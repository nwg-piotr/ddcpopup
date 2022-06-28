build:
	go build -o bin/ddcpopup .

install:
	cp bin/ddcpopup /usr/bin

uninstall:
	rm /usr/bin/ddcpopup

run:
	go run .