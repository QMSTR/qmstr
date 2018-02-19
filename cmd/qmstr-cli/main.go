//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto
package main

import (
	"github.com/QMSTR/qmstr/pkg/cli"
)

func main() {
	cli.Execute()
}
