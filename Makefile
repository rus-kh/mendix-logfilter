all: clean mendix-logfilter

GOPATH = $(CURDIR)/go

clean:
	rm -f mendix-logfilter
	rm -f glide.tar.gz
	rm -f glide
	rm -rf go/
	rm -rf vendor/

glide.tar.gz:
	curl -sL -o glide.tar.gz https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-linux-amd64.tar.gz

glide: glide.tar.gz
	tar --strip=1 -zxvf glide.tar.gz -- linux-amd64/glide
	mkdir $(GOPATH) vendor
	ln -sf ../vendor $(GOPATH)/src
	GOPATH=$(GOPATH) ./glide install

mendix-logfilter: glide
	GOPATH=$(GOPATH) go build -ldflags="-s -w" -o $@ ./main.go
