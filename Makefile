PROTO_PYTHON_FILES := $(shell find python/ -type f -name '*_pb2*.py' -printf '%p ')
PYTHON_FILES := $(filter-out $(PROTO_PYTHON_FILES), $(shell find python/ -type f -name '*.py' -printf '%p '))
GO_PKGS := $(shell go list ./... | grep -v /vendor)
GO_BIN := $(GOPATH)/bin
GOMETALINTER := $(GO_BIN)/gometalinter
GODEP := $(GO_BIN)/dep
QMSTR_GO_BINARIES := qmstr-wrapper qmstr qmstr-master qmstr-cli

.PHONY: all
all: $(QMSTR_GO_BINARIES)

generate: go_proto python_proto

go_proto:
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

python_proto:
	@mkdir python/pyqmstr/service || true
	python -m grpc_tools.protoc -Iproto --python_out=./python/pyqmstr/service --grpc_python_out=./python/pyqmstr/service proto/*.proto

.PHONY: clean
clean:
	@rm python/pyqmstr/service/*_pb2*.py || true
	@rm pkg/service/*.pb.go || true
	@rm $(QMSTR_GO_BINARIES) || true

.PHONY: checkpep8
checkpep8: $(PYTHON_FILES)
	@autopep8 --diff $^

.PHONY: checkpep8
autopep8: $(PYTHON_FILES)
	@autopep8 -i $^

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

$(QMSTR_GO_BINARIES): godep go_proto gotest
	go build github.com/QMSTR/qmstr/cmd/$@

.PHONY: container
container: ci/Dockerfile
	docker build --no-cache -f ci/Dockerfile -t qmstr/master --target master .

.PHONY: devcontainer
devcontainer: ci/Dockerfile
	docker build --no-cache -f ci/Dockerfile -t qmstr/dev --target dev .
