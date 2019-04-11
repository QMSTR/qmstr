GO_MODULE_PKGS := $(shell go list ./... | grep /module | grep -v /vendor)
GO_PKGS := $(shell go list ./... | grep -v /module | grep -v /vendor)

# those are intended to be a recursively expanded variables
GO_SRCS = $(shell find cmd -name '*.go') $(shell find pkg -name '*.go') 
FILTER = $(foreach v,$(2),$(if $(findstring $(1),$(v)),$(v),))
GO_MODULE_SRCS = $(call FILTER,module,$(GO_SRCS))

GO_PATH := $(shell go env GOPATH)
GO_BIN := $(GO_PATH)/bin
GOMETALINTER := $(GO_BIN)/gometalinter
GODEP := $(GO_BIN)/dep

PROTOC_GEN_GO := $(GO_BIN)/protoc-gen-go
PROTOC_GEN_GO_SRC := vendor/github.com/golang/protobuf/protoc-gen-go

OUTDIR := out/
QMSTR_GO_ANALYZERS := $(foreach ana, $(shell ls cmd/modules/analyzers), ${OUTDIR}analyzers/$(ana))
QMSTR_GO_REPORTERS := $(foreach rep, $(shell ls cmd/modules/reporters), ${OUTDIR}reporters/$(rep))
QMSTR_GO_BUILDERS := $(foreach builder, qmstr-wrapper, ${OUTDIR}$(builder))
QMSTR_CLIENT_BINARIES := $(foreach cli, qmstrctl qmstr, ${OUTDIR}$(cli)) $(QMSTR_GO_BUILDERS)
QMSTR_MASTER := $(foreach bin, qmstr-master, ${OUTDIR}$(bin))
QMSTR_SERVER_BINARIES := $(QMSTR_MASTER) $(QMSTR_GO_ANALYZERS) $(QMSTR_GO_REPORTERS) $(QMSTR_PYTHON_MODULES)
QMSTR_GO_BINARIES := $(QMSTR_MASTER) $(QMSTR_CLIENT_BINARIES)
QMSTR_GO_MODULES := $(QMSTR_GO_ANALYZERS) $(QMSTR_GO_REPORTERS)

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
	@rm -r out || true
	@rm -fr vendor || true
	@rm .go_module_test .go_qmstr_test || true

.PHONY: cleanall
cleanall: clean
	@docker rmi ${CONTAINER_TAG_DEV} ${CONTAINER_TAG_MASTER} ${CONTAINER_TAG_RUNTIME} || true
	@docker image prune --all --force --filter=label=org.qmstr.image

.go_module_test: $(GO_MODULE_SRCS)
	go test $(GO_MODULE_PKGS)
	@touch .go_module_test

.go_qmstr_test: $(GO_SRCS)
	go test $(GO_PKGS)
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

$(GODEP):
	go get -u -v github.com/golang/dep/cmd/dep

Gopkg.lock: Gopkg.toml
	${GO_BIN}/dep ensure --no-vendor

vendor: $(GODEP) Gopkg.lock
	${GO_BIN}/dep ensure --vendor-only

$(QMSTR_GO_BINARIES): vendor $(GO_SRCS) .go_qmstr_test 
	go build -o $@ github.com/QMSTR/qmstr/cmd/$(subst $(OUTDIR),,$@)

$(QMSTR_GO_MODULES): vendor $(GO_SRCS) .go_module_test
	go build -o $@ github.com/QMSTR/qmstr/cmd/modules/$(subst $(OUTDIR),,$@)

.PHONY: container master
master container: ci/Dockerfile
	docker build -f ci/Dockerfile -t ${CONTAINER_TAG_MASTER} --target master $(DOCKER_PROXY) .

.PHONY: devcontainer
devcontainer: container
	docker build -f ci/Dockerfile -t ${CONTAINER_TAG_DEV} --target dev $(DOCKER_PROXY) .

.PHONY: democontainer
democontainer: container
	docker build -f ci/Dockerfile -t ${CONTAINER_TAG_RUNTIME} --target runtime $(DOCKER_PROXY) .
	docker build -f ci/Dockerfile -t ${CONTAINER_TAG_BUILDER} --target builder $(DOCKER_PROXY) .

.PHONY: ratelimage
ratelimage:
	docker build -f ci/Dockerfile -t qmstr-ratel --target web $(DOCKER_PROXY) .

install_qmstr_server: $(QMSTR_SERVER_BINARIES)
	cp $^ /usr/local/bin

install_qmstr_client: $(QMSTR_CLIENT_BINARIES)
	cp $^ /usr/local/bin

install_qmstr_all: install_qmstr_client install_qmstr_server

install_qmstr_client_gopath: $(QMSTR_CLIENT_BINARIES)
	cp $^ ${GO_PATH}/bin/

pkg/service/%.pb.go: $(PROTOC_GEN_GO) proto/%.proto
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

# Generate doc website statics files
.PHONY: docs
docs:
	cd doc && hugo -d documentation
	cd doc && tar -cjvf qmstr-doc.tar.bz2 documentation
	rm -rf doc/documentation
