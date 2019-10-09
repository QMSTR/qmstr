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
	Input      []string
	Output     string
	OutputType string
	builder.GeneralBuilder
}

func NewObjcopyBuilder(workDir string, logger *log.Logger, debug bool) *ObjcopyBuilder {
	return &ObjcopyBuilder{[]string{}, workDir, []string{}, []string{}, "", "", builder.NewGeneralBuilder(logger, debug)}
}

func (o *ObjcopyBuilder) GetPrefix() (string, error) {
	return "", errors.New("objcopy not prefixed")
}

func (o *ObjcopyBuilder) GetName() string {
	return "objcopy builder"
}

func (o *ObjcopyBuilder) Analyze(commandline []string) ([]*service.FileNode, error) {
	if o.Debug {
		o.Logger.Printf("%s parsing commandline %v", o.GetName(), commandline)
	}

	err := o.processFlags(commandline[1:])
	if err != nil {
		return nil, err
	}

	dependencies := []*service.FileNode{}
	outputTarget := builder.NewFileNode(common.BuildCleanPath(o.Workdir, o.Output, false), false)
	for _, input := range o.Input {
		inputTarget := builder.NewFileNode(common.BuildCleanPath(o.Workdir, input, false), true)
		dependencies = append(dependencies, inputTarget)
		o.Logger.Printf("%s copying from %s:%s to %s", o.GetName(), inputTarget.GetPath(), inputTarget.FileData.GetHash(), o.Output)
	}
	outputTarget.DerivedFrom = dependencies
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
		case "-S", "-p", "-D", "-U", "-g", "-w", "-x", "-X", "-M", "--only-keep-debug", "--compress-debug-sections":
			cleanIdx = append(cleanIdx, idx)
			continue
		}
	}

	o.Args = builder.CleanCmd(args, cleanIdx, o.Debug, o.Logger)

	objCpFlags := pflag.NewFlagSet("objcopy", pflag.ContinueOnError)
	objCpFlags.StringP("output-target", "O", undef, "output-target")
	objCpFlags.BoolP("strip-debug", "g", false, "do not copy debugging symbols or sections from the source file")
	objCpFlags.String("add-gnu-debuglink", undef, "create a .gnu_debuglink section")
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
	o.Input = append(o.Input, o.Targets[0])
	if len(o.Targets) == 1 {
		if outputType, err := objCpFlags.GetString("output-target"); err == nil && outputType != undef {
			switch outputType {
			case "binary":
				o.Output = strings.TrimSuffix(o.Input[0], filepath.Ext(o.Input[0])) + ".bin"
			default:
				return errors.New("Output format not implemented")
			}
		}
		if targetType, err := objCpFlags.GetString("target"); err == nil && targetType != undef {
			switch targetType {
			case "binary":
				o.Output = strings.TrimSuffix(o.Input[0], filepath.Ext(o.Input[0])) + ".bin"
			default:
				return errors.New("Output format not implemented")
			}
		}
		if stripBool, err := objCpFlags.GetBool("strip-debug"); err == nil && stripBool {
			o.Output = o.Input[0]
		}
		if debugFile, err := objCpFlags.GetString("add-gnu-debuglink"); err == nil && debugFile != undef {
			o.Output = o.Input[0]
			o.Input = append(o.Input, debugFile)
		}
		return nil
	}
	o.Output = o.Targets[1]
	return nil
}
