PROTO_PYTHON_FILES = $(shell find python/ -type f -name '*_pb2*.py' -printf '%p ')
PYTHON_FILES = $(filter-out $(PROTO_PYTHON_FILES), $(shell find python/ -type f -name '*.py' -printf '%p '))
GO_PKGS = $(shell go list ./... | grep -v /vendor)

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