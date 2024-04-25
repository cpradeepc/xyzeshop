# variables
GOBUILD=go build  
GOTEST=go test 


# multi make  commands in all make command
all:  stop  cls build
	# ./xyzeshop.exe
build:
	$(GOBUILD) -v .
cls:
	rm -f ./xyzeshop
stop:
	 pkil -15 xyzeshop || true
	# taskkill /IM xyzeshop.exe /F || true
run:
	./xyzeshop
# test:
# 	cd helper && $(GOTEST) -v .