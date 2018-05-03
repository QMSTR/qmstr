generate: go_proto python_proto

go_proto:
	protoc -I proto --go_out=plugins=grpc:pkg/service proto/*.proto

python_proto:
	mkdir python/service || true
	python -m grpc_tools.protoc -Iproto --python_out=./python/service --grpc_python_out=./python/service proto/*.proto

clean:
	rm python/service/*_pb2*.py
	rm pkg/service/*.pb.go
