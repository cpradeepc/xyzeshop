# variables
GOBUILD=go build  
GOTEST=go test 


# multi make  commands in all make command
all: cls stop build_server 
	./xyzeshop
build_server:
	$(GOBUILD) -v .
cls:
	rm -f ./xyzeshop
stop:
	pkil xyzeshop || true 
test:
	cd helper && $(GOTEST) -v .