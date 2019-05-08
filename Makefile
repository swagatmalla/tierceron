GOPATH=/home/mjrieke/workspace/go/:$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)

api:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/apiRouter webapi/server/server.go
config:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/vaultconfig vaultconfig/vaultconfig.go
seed:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/vaultinit vaultinit/init.go
gen:
	protoc --proto_path=. --twirp_out=. --go_out=. rpc/apinator/service.proto
