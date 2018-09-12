package builder

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	pb "github.com/QMSTR/qmstr/pkg/service"
)

type Builder interface {
	Analyze(commandline []string) (*pb.BuildMessage, error)
	GetName() string
	GetPrefix() (string, error)
}

type GeneralBuilder struct {
	Logger *log.Logger
	Debug  bool
}

func NewFileNode(path string, fileType string) *pb.FileNode {
	filename := filepath.Base(path)
	hash, err := hash(path)
	broken := false
	if err != nil {
		hash = "nohash" + path
		broken = true
	}
	return &pb.FileNode{NodeType: pb.NodeTypeFileNode, Name: filename, Type: fileType, Path: path, Hash: hash, Broken: broken}
}

func hash(fileName string) (string, error) {
	h := sha1.New()
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := f.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return "", err
		}
		h.Write(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
