# Versions
GRPCIO_VERSION := 1.15.0

# Common stuff
OUTDIR := out/

PROTO_FILES := $(shell ls ./proto/*.proto)
ALL_PYTHON_FILES := $(shell find . -path ./venv -prune -o -name *.py -print)
PROTO_PYTHON_FILES := $(patsubst ./proto/%,./lib/pyqmstr/qmstr/service/%,$(PROTO_FILES:.proto=_pb2.py)) $(patsubst ./proto/%,./lib/pyqmstr/qmstr/service/%,$(PROTO_FILES:.proto=_pb2_grpc.py))
PYTHON_FILES := $(filter-out $(PROTO_PYTHON_FILES), $(ALL_PYTHON_FILES))

GO_MODULE_PKGS := $(shell go list ./... | grep /module | grep -v /vendor)
GO_PKGS := $(shell go list ./... | grep -v /module | grep -v /vendor)

# those are intended to be a recursively expanded variables
GO_SRCS = $(shell find modules -name '*.go') $(shell find pkg -name '*.go') 
FILTER = $(foreach v,$(2),$(if $(findstring $(1),$(v)),$(v),))
GO_MODULE_SRCS = $(call FILTER,module,$(GO_SRCS))

GO_PATH := $(shell go env GOPATH)
GO_BIN := $(GO_PATH)/bin
GOMETALINTER := $(GO_BIN)/gometalinter

PROTOC_GEN_GO := $(GO_BIN)/protoc-gen-go

QMSTR_ANALYZERS := $(foreach ana, $(shell ls modules/analyzers), $(OUTDIR)analyzers/$(ana))
QMSTR_REPORTERS := $(foreach rep, $(shell ls modules/reporters), $(OUTDIR)reporters/$(rep))
QMSTR_BUILDERS := $(foreach builder, $(shell ls modules/builders), $(OUTDIR)builders/$(builder))

QMSTR_CLIENT_BINARIES := $(foreach cli, qmstrctl qmstr, ${OUTDIR}$(cli)) $(QMSTR_BUILDERS)
QMSTR_MASTER := $(foreach bin, qmstr-master, ${OUTDIR}$(bin))
QMSTR_SERVER_BINARIES := $(QMSTR_MASTER) $(QMSTR_REPORTERS) $(QMSTR_ANALYZERS)

CONTAINER_TAG_DEV := qmstr/dev
CONTAINER_TAG_MASTER := qmstr/master
CONTAINER_TAG_RUNTIME := qmstr/runtime
CONTAINER_TAG_BUILDER := qmstr/master_build

ifdef http_proxy
	DOCKER_PROXY = --build-arg http_proxy=$(http_proxy)
endif

ifdef https_proxy
	DOCKER_PROXY += --build-arg https_proxy=$(https_proxy)
endif


.PHONY: all
all: install_qmstr_client_gopath democontainer

.PHONY: clean
clean:
	git clean -ffxd

.PHONY: cleanall
cleanall: clean
	@docker rmi ${CONTAINER_TAG_DEV} ${CONTAINER_TAG_MASTER} ${CONTAINER_TAG_RUNTIME} || true
	@docker image prune --all --force --filter=label=org.qmstr.image

.go_module_test: $(GO_MODULE_SRCS)
	go test ./modules/...
	@touch .go_module_test

.go_qmstr_test: $(GO_SRCS)
	go test ./cmd/... ./pkg/...
	@touch .go_qmstr_test

.PHONY: gotest
gotest: .go_qmstr_test .go_module_test

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: golint
golint:	$(GOMETALINTER)
	gometalinter ./... --vendor

.PHONY: govet
govet: gotest
	go tool vet cmd pkg

install_qmstr_server: $(QMSTR_SERVER_BINARIES)
	cp $^ /usr/local/bin

install_qmstr_client: $(QMSTR_CLIENT_BINARIES)
	cp $^ /usr/local/bin

install_qmstr_all: install_qmstr_client install_qmstr_server

install_qmstr_client_gopath: $(QMSTR_CLIENT_BINARIES)
	cp $^ ${GO_PATH}/bin/

pkg/service/%.pb.go: $(PROTOC_GEN_GO) proto/%.proto
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

# Python related stuff
venv: venv/bin/activate
venv/bin/activate: requirements.txt
	test -d venv || virtualenv -p python3 venv
	venv/bin/pip install -Ur requirements.txt
	touch venv/bin/activate

requirements.txt:
	@echo grpcio==$(GRPCIO_VERSION) >> requirements.txt
	@echo grpcio-tools==$(GRPCIO_VERSION) >> requirements.txt

venv/bin/pex: venv
	@venv/bin/pip install pex

venv/bin/autopep8: venv
	@venv/bin/pip install autopep8

.PHONY: checkpep8
checkpep8: $(PYTHON_FILES) venv
	venv/bin/autopep8 --diff $(filter-out venv, $^)

.PHONY: autopep8
autopep8: $(PYTHON_FILES) venv
	venv/bin/autopep8 -i $(filter-out venv, $^)


# include all makefiles
include $(shell find . -name "Makefile.common")