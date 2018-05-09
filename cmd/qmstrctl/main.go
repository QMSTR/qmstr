//go:generate protoc -I ../../proto --go_out=plugins=grpc:../../pkg/service ../../proto/datamodel.proto ../../proto/analyzerservice.proto ../../proto/buildservice.proto ../../proto/controlservice.proto  ../../proto/reportservice.proto
package main

import (
	"github.com/QMSTR/qmstr/pkg/cli"
)

func main() {
	cli.Execute()
}
