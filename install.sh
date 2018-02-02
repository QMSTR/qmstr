#!/bin/bash
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/golang/dep/cmd/dep

go get github.com/QMSTR/qmstr
(cd $GOPATH/src/github.com/QMSTR/qmstr; dep ensure)
go generate github.com/QMSTR/qmstr/cmd/qmstr-master

go install github.com/QMSTR/qmstr/cmd/qmstr-master
go install github.com/QMSTR/qmstr/cmd/qmstr-wrapper
go install github.com/QMSTR/qmstr/cmd/qmstr-client
