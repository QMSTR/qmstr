PROTO_PYTHON_FILES = $(shell find python/ -type f -name '*_pb2*.py' -printf '%p ')
PYTHON_FILES = $(filter-out $(PROTO_PYTHON_FILES), $(shell find python/ -type f -name '*.py' -printf '%p '))

generate: go_proto python_proto

go_proto:
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

python_proto:
	mkdir python/pyqmstr/service || true
	python -m grpc_tools.protoc -Iproto --python_out=./python/pyqmstr/service --grpc_python_out=./python/pyqmstr/service proto/*.proto

clean:
	rm python/pyqmstr/service/*_pb2*.py
	rm pkg/service/*.pb.go

checkpep8: $(PYTHON_FILES)
	autopep8 --diff $^
