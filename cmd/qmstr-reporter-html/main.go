//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto
package main

import (
	"github.com/endocode/qmstr/pkg/reporter/htmlreporter"
)

func main() {
	htmlreporter.ConnectToMaster()
	htmlreporter.Temp()
	defer htmlreporter.DisconnectFromMaster()
}
