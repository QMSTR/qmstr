package logging

import (
	"context"
	"errors"

	pb "github.com/QMSTR/go-qmstr/service"
)

// RemoteLogWriter can be used as logger sink that sends the log messages to the server
type RemoteLogWriter struct {
	bsc pb.ControlServiceClient
}

func NewRemoteLogWriter(bsc pb.ControlServiceClient) *RemoteLogWriter {
	return &RemoteLogWriter{bsc: bsc}
}

func (rlw *RemoteLogWriter) Write(p []byte) (int, error) {
	logmsg := pb.LogMessage{Msg: p}
	r, err := rlw.bsc.Log(context.Background(), &logmsg)
	if err != nil {
		return 0, err
	}
	if !r.Success {
		return 0, errors.New("server failed to process log message")
	}
	return len(p), nil
}
