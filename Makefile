# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=FileToMQSender
BINARY_UNIX=$(BINARY_NAME)_linux

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

build-release-pj:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_UNIX) -v
	cp pj_config.toml release/config.toml
	cp start.sh release/start.sh
	cp kill.sh release/kill.sh
	cp status.sh release/status.sh
	sshpass -p "pass"  scp -r release/*  root@10.11.2.60:/opt/app/huadi/op_go_reader/release_pj/


build-release-deppon:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_UNIX) -v
	cp deppon_config.toml release/config.toml
	cp start.sh release/start.sh
	cp kill.sh release/kill.sh
	cp status.sh release/status.sh

	sshpass -p "pass"  scp -r release/*  root@10.11.2.60:/opt/app/huadi/op_go_reader/release_deppon/

build-release-sto:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_UNIX) -v
	cp sto_config.toml release/config.toml
	cp start.sh release/start.sh
	cp kill.sh release/kill.sh
	cp status.sh release/status.sh
	tar -czvf release_sto.tar.gz release