PROTO_PYTHON_FILES := $(shell find python/ -type f -name '*_pb2*.py' -printf '%p ')
PYTHON_FILES := $(filter-out $(PROTO_PYTHON_FILES), $(shell find python/ -type f -name '*.py' -printf '%p '))
GO_PKGS := $(shell go list ./... | grep -v /vendor)
GO_BIN := $(GOPATH)/bin
GOMETALINTER := $(GO_BIN)/gometalinter
GODEP := $(GO_BIN)/dep
QMSTR_GO_BINARIES := qmstr-wrapper qmstr qmstr-master qmstr-cli

GRPCIOTOOLS_VERSION := 1.11.0

.PHONY: all
all: $(QMSTR_GO_BINARIES)

generate: go_proto python_proto

venv: venv/bin/activate
venv/bin/activate: requirements.txt 
	test -d venv || virtualenv venv
	venv/bin/pip install -Ur requirements.txt
	touch venv/bin/activate

requirements.txt:
	echo grpcio-tools==$(GRPCIOTOOLS_VERSION) >> requirements.txt
	echo pex >> requirements.txt

go_proto:
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

python_proto: venv
	@mkdir python/pyqmstr/service || true
	venv/bin/python -m grpc_tools.protoc -Iproto --python_out=./python/pyqmstr/service --grpc_python_out=./python/pyqmstr/service proto/*.proto

.PHONY: clean
clean:
	@rm python/pyqmstr/service/*_pb2*.py || true
	@rm pkg/service/*.pb.go || true
	@rm $(QMSTR_GO_BINARIES) || true
	@rm -fr venv
	@rm requirements.txt

.PHONY: checkpep8
checkpep8: $(PYTHON_FILES)
	@autopep8 --diff $^

.PHONY: autopep8
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
devcontainer: container
	docker build -f ci/Dockerfile -t qmstr/dev --target dev .
