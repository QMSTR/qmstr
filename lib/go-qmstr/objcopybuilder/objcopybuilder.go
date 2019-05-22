package objcopybuilder

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/QMSTR/qmstr/lib/go-qmstr/builder"
	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"

	"github.com/spf13/pflag"
)

const undef = "undef"

type ObjcopyBuilder struct {
	Args       []string
	Workdir    string
	Targets    []string
	Input      string
	Output     string
	OutputType string
	builder.GeneralBuilder
}

func NewObjcopyBuilder(workDir string, logger *log.Logger, debug bool) *ObjcopyBuilder {
	return &ObjcopyBuilder{[]string{}, workDir, []string{}, "", "", "", builder.NewGeneralBuilder(logger, debug)}
}

func (o *ObjcopyBuilder) GetPrefix() (string, error) {
	return "", errors.New("objcopy not prefixed")
}

func (o *ObjcopyBuilder) GetName() string {
	return "objcopy builder"
}

func (o *ObjcopyBuilder) Analyze(commandline []string) ([]*service.FileNode, error) {
	o.Logger.Printf("Objcopy copying binary file")

	if o.Debug {
		o.Logger.Printf("Parsing commandline %v", commandline)
	}

	err := o.processFlags(commandline[1:])
	if err != nil {
		return nil, err
	}

	outputTarget := builder.NewFileNode(common.BuildCleanPath(o.Workdir, o.Output, false), service.FileNode_TARGET, false)
	inputTarget := builder.NewFileNode(common.BuildCleanPath(o.Workdir, o.Input, false), service.FileNode_TARGET, true)
	outputTarget.DerivedFrom = []*service.FileNode{inputTarget}

	return []*service.FileNode{outputTarget}, nil
}

func (o *ObjcopyBuilder) processFlags(args []string) error {
	if o.Debug {
		o.Logger.Printf("Parsing arguments: %v", args)
	}

	cleanIdx := []int{}
	for idx, arg := range args {
		switch arg {
		case "-R", "-j", "-K", "-N", "-L", "-G", "-W", "-b":
			cleanIdx = append(cleanIdx, idx, idx+1)
			continue
		case "-S", "-p", "-D", "-U", "-g", "-w", "-x", "-X", "-M":
			cleanIdx = append(cleanIdx, idx)
			continue
		}
	}

	o.Args = builder.CleanCmd(args, cleanIdx, o.Debug, o.Logger)

	objCpFlags := pflag.NewFlagSet("objcopy", pflag.ContinueOnError)
	objCpFlags.StringP("output-target", "O", undef, "output-target")
	objCpFlags.StringP("input-target", "I", undef, "input-target")
	objCpFlags.StringP("target", "F", undef, "target")

	if o.Debug {
		o.Logger.Printf("Parsing cleaned commandline: %v", o.Args)
	}

	err := objCpFlags.Parse(o.Args)
	if err != nil {
		return fmt.Errorf("Unrecoverable commandline parsing error: %v", err)
	}

	o.Targets = objCpFlags.Args()

	if len(o.Targets) <= 0 {
		return builder.ErrNoTargetsProvided
	}
	o.Input = o.Targets[0]
	if len(o.Targets) == 1 {
		if outputType, err := objCpFlags.GetString("output-target"); err == nil && outputType != undef {
			switch outputType {
			case "binary":
				o.Output = strings.TrimSuffix(o.Input, filepath.Ext(o.Input)) + ".bin"
			default:
				return errors.New("Output format not implemented")
			}
		}
		if targetType, err := objCpFlags.GetString("target"); err == nil && targetType != undef {
			switch targetType {
			case "binary":
				o.Output = strings.TrimSuffix(o.Input, filepath.Ext(o.Input)) + ".bin"
			default:
				return errors.New("Output format not implemented")
			}
		}
		return nil
	}
	o.Output = o.Targets[1]
	return nil
}
