package compiler

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	pb "github.com/QMSTR/qmstr/pkg/service"
	"github.com/QMSTR/qmstr/pkg/wrapper"
	"github.com/spf13/pflag"
)

type mode int

const (
	Link mode = iota
	Preproc
	Compile
	Assemble
	PrintOnly
)
const undef = "undef"
const (
	linkedTrg = "linkedtarget"
	obj       = "objectfile"
	src       = "sourcefile"
	lib       = "library"
)

var (
	boolArgs = map[string]struct{}{
		"-w":               struct{}{},
		"-W":               struct{}{},
		"-O":               struct{}{},
		"-f":               struct{}{},
		"-C":               struct{}{},
		"-std":             struct{}{},
		"-nostdinc":        struct{}{},
		"-nostdlib":        struct{}{},
		"-print-file-name": struct{}{},
		"-M":               struct{}{},
		"-MG":              struct{}{},
		"-MM":              struct{}{},
		"-MD":              struct{}{},
		"-MMD":             struct{}{},
		"-MP":              struct{}{},
		"-m":               struct{}{},
		"-v":               struct{}{},
		"-g":               struct{}{},
		"-pg":              struct{}{},
		"-P":               struct{}{},
		"-pipe":            struct{}{},
		"-rdynamic":        struct{}{},
		"-shared":          struct{}{},
		"-static":          struct{}{},
		"-dynamiclib":      struct{}{},
		"--version":        struct{}{},
	}

	stringArgs = map[string]struct{}{
		"-D":                     struct{}{},
		"-U":                     struct{}{},
		"-x":                     struct{}{},
		"-MF":                    struct{}{},
		"-MT":                    struct{}{},
		"-MQ":                    struct{}{},
		"-install_name":          struct{}{},
		"-compatibility_version": struct{}{},
		"-current_version":       struct{}{},
	}

	stringArgsRE = "\\s+\\S+={0,1}\\S*\\s"

	fixPosixArgs = map[string]struct{}{
		"-isystem": struct{}{},
		"-include": struct{}{},
	}
)

type GccCompiler struct {
	Mode     mode
	Input    []string
	Output   []string
	WorkDir  string
	LinkLibs []string
	LibPath  []string
	Args     []string
	GeneralCompiler
}

func NewGccCompiler(workDir string, logger *log.Logger, debug bool) *GccCompiler {
	return &GccCompiler{Link, []string{}, []string{}, workDir, []string{}, []string{}, []string{}, GeneralCompiler{logger, debug}}
}

func (g *GccCompiler) Analyze(commandline []string) (*pb.BuildMessage, error) {
	if g.debug {
		g.logger.Printf("Parsing commandline %v", commandline)
	}
	g.parseCommandLine(commandline[1:])

	switch g.Mode {
	case Link:
		g.logger.Printf("gcc linking")
		linkedTarget := NewFileNode(wrapper.BuildCleanPath(g.WorkDir, g.Output[0], false), linkedTrg)
		fileNodes := []*pb.FileNode{linkedTarget}
		dependencies := []*pb.FileNode{}
		for _, inFile := range g.Input {
			inputFileNode := &pb.FileNode{}
			ext := filepath.Ext(inFile)
			if ext == ".o" {
				inputFileNode = NewFileNodeDerivedFrom(wrapper.BuildCleanPath(g.WorkDir, inFile, false), obj, fileNodes)
			} else if ext == ".c" {
				inputFileNode = NewFileNodeDerivedFrom(wrapper.BuildCleanPath(g.WorkDir, inFile, false), src, fileNodes)
			} else {
				inputFileNode = NewFileNodeDerivedFrom(wrapper.BuildCleanPath(g.WorkDir, inFile, false), lib, fileNodes)
			}
			dependencies = append(dependencies, inputFileNode)
			fileNodes = append(fileNodes, inputFileNode)
		}

		actualLibs, err := wrapper.FindActualLibraries(g.LinkLibs, g.LibPath)
		if err != nil {
			g.logger.Fatalf("Failed to collect dependencies: %v", err)
		}
		for _, lib := range actualLibs {
			linkLib := NewFileNodeDerivedFrom(lib, lib, dependencies)
			dependencies = append(dependencies, linkLib)
			fileNodes = append(fileNodes, linkLib)
		}
		return &pb.BuildMessage{fileNodes}, nil
	case Assemble:
		g.logger.Printf("gcc assembling - skipping link")
		fileNodes := []*pb.FileNode{}
		if g.debug {
			g.logger.Printf("This is our input %v", g.Input)
			g.logger.Printf("This is our output %v", g.Output)
		}
		for idx, inFile := range g.Input {
			if g.debug {
				g.logger.Printf("This is the source file %s indexed %d", inFile, idx)
			}
			targetFile := NewFileNode(wrapper.BuildCleanPath(g.WorkDir, g.Output[idx], false), obj)
			fileNodes = append(fileNodes, targetFile)
			sourceFile := NewFileNodeDerivedFrom(wrapper.BuildCleanPath(g.WorkDir, inFile, false), src, fileNodes)
			fileNodes = append(fileNodes, sourceFile)
		}
		return &pb.BuildMessage{fileNodes}, nil
	default:
		return nil, errors.New("Mode not implemented")
	}
}

func (g *GccCompiler) cleanCmdLine(args []string) {
	clearIdxSet := map[int]struct{}{}
	for idx, arg := range args {

		if g.debug {
			g.logger.Printf("%d - %s", idx, arg)
		}

		// index string flags
		if idx < len(args)-1 {
			for key := range stringArgs {
				if g.debug {
					g.logger.Printf("Find %s string arg in %s with %s", key, fmt.Sprintf("%s %s ", arg, args[idx+1]), fmt.Sprintf("%s%s", key, stringArgsRE))
				}
				re := regexp.MustCompile(fmt.Sprintf("%s%s", key, stringArgsRE))
				if re.MatchString(fmt.Sprintf("%s %s ", arg, args[idx+1])) {
					if g.debug {
						g.logger.Printf("Found %v string arg", args[idx:idx+1])
					}
					clearIdxSet[idx] = struct{}{}
					clearIdxSet[idx+1] = struct{}{}
				}
				if strings.HasPrefix(arg, key) {
					clearIdxSet[idx] = struct{}{}
				}
			}
		}

		// index bool flags
		for key := range boolArgs {
			if strings.HasPrefix(arg, key) {
				clearIdxSet[idx] = struct{}{}
			}
		}

		// fix long arguments to pass through pflags
		for key := range fixPosixArgs {
			if key == arg {
				args[idx] = fmt.Sprintf("-%s", arg)
			}
		}
	}

	clear := []int{}
	for k := range clearIdxSet {
		clear = append(clear, k)
	}
	sort.Sort(sort.IntSlice(clear))

	if g.debug {
		g.logger.Printf("To be cleaned %v", clear)
	}
	initialArgsSize := len(args)
	for _, idx := range clear {
		if g.debug {
			g.logger.Printf("Clearing %d", idx)
		}
		offset := initialArgsSize - len(args)
		offsetIdx := idx - offset
		if g.debug {
			g.logger.Printf("Actually clearing %d", offsetIdx)
		}
		if initialArgsSize-1 == idx {
			if g.debug {
				g.logger.Printf("Cut last arg")
			}
			args = args[:offsetIdx]
		} else {
			args = append(args[:offsetIdx], args[offsetIdx+1:]...)
		}
		if g.debug {
			g.logger.Printf("new slice is %v", args)
		}
	}
	g.Args = args
}

func (g *GccCompiler) parseCommandLine(args []string) {
	if g.debug {
		g.logger.Printf("Parsing arguments: %v", args)
	}

	// remove all flags we don't care about but that would break parsing
	g.cleanCmdLine(args)

	gccFlags := pflag.NewFlagSet("gcc", pflag.ContinueOnError)
	gccFlags.BoolP("assemble", "c", false, "do not link")
	gccFlags.BoolP("compile", "S", false, "do not assemble")
	gccFlags.BoolP("preprocess", "E", false, "do not compile")
	gccFlags.StringP("output", "o", undef, "output")
	gccFlags.StringSliceP("includepath", "I", []string{}, "include path")
	gccFlags.String("isystem", undef, "system include path")
	gccFlags.String("include", undef, "include header file")
	gccFlags.StringSliceVarP(&g.LinkLibs, "linklib", "l", []string{}, "link libs")
	gccFlags.StringSliceVarP(&g.LibPath, "linklibdir", "L", []string{}, "search dir for link libs")

	if g.debug {
		g.logger.Printf("Parsing cleaned commandline: %v", g.Args)
	}
	err := gccFlags.Parse(g.Args)
	if err != nil {
		g.logger.Fatalf("Unrecoverable commandline parsing error: %s", err)
	}

	g.Input = gccFlags.Args()

	if ok, err := gccFlags.GetBool("assemble"); ok && err == nil {
		g.Mode = Assemble
	}
	if ok, err := gccFlags.GetBool("compile"); ok && err == nil {
		g.Mode = Compile
	}
	if ok, err := gccFlags.GetBool("preprocess"); ok && err == nil {
		g.Mode = Preproc
	}
	if g.debug {
		g.logger.Printf("Mode set to: %v", g.Mode)
	}

	if output, err := gccFlags.GetString("output"); err == nil && output != undef {
		g.Output = []string{output}
	} else {
		// no output defined
		switch g.Mode {
		case Link:
			if len(g.Input) == 0 {
				// No input no output
				g.Mode = PrintOnly
				return
			}
			g.Output = []string{"a.out"}
		case Assemble:
			for _, input := range g.Input {
				objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".o"
				g.Output = append(g.Output, objectname)
			}
		}
	}
}
