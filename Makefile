build:
	go build -o bin/ddcpopup .

install:
	cp bin/ddcpopup /usr/bin
	mkdir -p /usr/share/ddcpopup
	cp -r icons /usr/share/ddcpopup

uninstall:
	rm -r /usr/share/ddcpopup
	rm /usr/bin/ddcpopup

run:
	go run .