//go:generate protoc -I ../../pkg/service --go_out=plugins=grpc:../../pkg/service ../../pkg/service/controlservice.proto
package main

import (
	"github.com/QMSTR/qmstr/pkg/cli"
)

func main() {
	cli.Execute()
}
