#!/bin/bash
go get -u github.com/golang/protobuf/protoc-gen-go

(go get github.com/QMSTR/qmstr/cmd/qmstr-master || go generate github.com/QMSTR/qmstr/cmd/qmstr-master; go get github.com/QMSTR/qmstr/cmd/qmstr-master)

go get github.com/QMSTR/qmstr/cmd/qmstr-wrapper
go get github.com/QMSTR/qmstr/cmd/qmstr-client