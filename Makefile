PROTO_PYTHON_FILES = $(shell find python/ -type f -name '*_pb2*.py' -printf '%p ')
PYTHON_FILES = $(filter-out $(PROTO_PYTHON_FILES), $(shell find python/ -type f -name '*.py' -printf '%p '))
GO_PKGS = $(shell go list ./... | grep -v /vendor)
GO_BIN = $(GOPATH)/bin
GOMETALINTER = $(GO_BIN)/gometalinter
GODEP = $(GO_BIN)/dep
QMSTR_MASTER = $(GO_BIN)/qmstr-master
QMSTR_WRAPPER = $(GO_BIN)/qmstr-wrapper
QMSTR_EASYMODE = $(GO_BIN)/qmstr
QMSTR_CLI= $(GO_BIN)/qmstr-cli

generate: go_proto python_proto

go_proto:
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

python_proto:
	@mkdir python/pyqmstr/service || true
	python -m grpc_tools.protoc -Iproto --python_out=./python/pyqmstr/service --grpc_python_out=./python/pyqmstr/service proto/*.proto

java_proto:
	@echo "Not yet implemented"

clean:
	@rm python/pyqmstr/service/*_pb2*.py
	@rm pkg/service/*.pb.go

checkpep8: $(PYTHON_FILES)
	@autopep8 --diff $^

.PHONY: gotest
gotest:
	go test $(GO_PKGS)

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: golint
golint:	$(GOMETALINTER)
	gometalinter ./... --vendor

$(GODEP):
	go get -u github.com/golang/dep

.PHONY: godep
godep: $(GODEP)
	dep ensure

qmstr-all-go: qmstr-master qmstr-wrapper qmstr-easymode qmstr-cli

qmstr-master: $(QMSTR_MASTER)

$(QMSTR_MASTER): godep go_proto gotest
	go install github.com/QMSTR/qmstr/cmd/qmstr-master

qmstr-wrapper: $(QMSTR_WRAPPER)

$(QMSTR_WRAPPER): godep go_proto gotest
	go install github.com/QMSTR/qmstr/cmd/qmstr-wrapper

qmstr-easymode: $(QMSTR_EASYMODE)

$(QMSTR_EASYMODE): godep go_proto gotest
	go install github.com/QMSTR/qmstr/cmd/qmstr

qmstr-cli: $(QMSTR_CLI)

$(QMSTR_CLI): godep go_proto gotest
	go install github.com/QMSTR/qmstr/cmd/qmstr-cli