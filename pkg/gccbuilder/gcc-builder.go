package gccbuilder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	pb "github.com/QMSTR/qmstr/pkg/service"
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
	src       = "sourcecode"
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
		"-pedantic":        struct{}{},
		"-print-":          struct{}{},
		"-pthread":         struct{}{},
		"-rdynamic":        struct{}{},
		"-shared":          struct{}{},
		"-static":          struct{}{},
		"-dynamiclib":      struct{}{},
		"-dumpversion":     struct{}{},
		"-dM":              struct{}{},
		"--version":        struct{}{},
		"-undef":           struct{}{},
		"-nostartfiles":    struct{}{},
		"-remap":           struct{}{},
		"-r":               struct{}{},
	}

	stringArgs = map[string]struct{}{
		"-D":                     struct{}{},
		"-Q":                     struct{}{},
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

type GccBuilder struct {
	Mode     mode
	Input    []string
	Output   []string
	WorkDir  string
	LinkLibs []string
	LibPath  []string
	Args     []string
	builder.GeneralBuilder
}

func NewGccBuilder(workDir string, logger *log.Logger, debug bool) *GccBuilder {
	return &GccBuilder{Link, []string{}, []string{}, workDir, []string{}, []string{}, []string{}, builder.GeneralBuilder{logger, debug}}
}

func (g *GccBuilder) GetName() string {
	return "GNU C compiler builder"
}

func (g *GccBuilder) Analyze(commandline []string) (*pb.BuildMessage, error) {
	if g.Debug {
		g.Logger.Printf("Parsing commandline %v", commandline)
	}
	g.parseCommandLine(commandline[1:])

	switch g.Mode {
	case Link:
		g.Logger.Printf("gcc linking")
		fileNodes := []*pb.FileNode{}
		linkedTarget := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[0], false), linkedTrg)
		dependencies := []*pb.FileNode{}
		for _, inFile := range g.Input {
			inputFileNode := &pb.FileNode{}
			ext := filepath.Ext(inFile)
			if ext == ".o" {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), obj)
			} else if ext == ".c" {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), src)
			} else {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), linkedTrg)
			}
			dependencies = append(dependencies, inputFileNode)
		}
		actualLibs, err := FindActualLibraries(g.LinkLibs, g.LibPath)
		if err != nil {
			g.Logger.Fatalf("Failed to collect dependencies: %v", err)
		}
		for _, actualLib := range actualLibs {
			linkLib := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, actualLib, false), linkedTrg)
			dependencies = append(dependencies, linkLib)
		}
		linkedTarget.DerivedFrom = dependencies
		fileNodes = append(fileNodes, linkedTarget)
		return &pb.BuildMessage{FileNodes: fileNodes}, nil
	case Assemble:
		g.Logger.Printf("gcc assembling - skipping link")
		fileNodes := []*pb.FileNode{}
		if g.Debug {
			g.Logger.Printf("This is our input %v", g.Input)
			g.Logger.Printf("This is our output %v", g.Output)
		}
		for idx, inFile := range g.Input {
			if g.Debug {
				g.Logger.Printf("This is the source file %s indexed %d", inFile, idx)
			}
			sourceFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), src)
			targetFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[idx], false), obj)
			targetFile.DerivedFrom = []*pb.FileNode{sourceFile}
			fileNodes = append(fileNodes, targetFile)
		}
		return &pb.BuildMessage{FileNodes: fileNodes}, nil
	default:
		return nil, errors.New("Mode not implemented")
	}
}

func (g *GccBuilder) cleanCmdLine(args []string) {
	clearIdxSet := map[int]struct{}{}
	for idx, arg := range args {

		if g.Debug {
			g.Logger.Printf("%d - %s", idx, arg)
		}

		// index string flags
		for key := range stringArgs {
			if idx < len(args)-1 {
				if g.Debug {
					g.Logger.Printf("Find %s string arg in %s with %s", key, fmt.Sprintf("%s %s ", arg, args[idx+1]), fmt.Sprintf("%s%s", key, stringArgsRE))
				}
				re := regexp.MustCompile(fmt.Sprintf("%s%s", key, stringArgsRE))
				if re.MatchString(fmt.Sprintf("%s %s ", arg, args[idx+1])) {
					if g.Debug {
						g.Logger.Printf("Found %v string arg", args[idx:idx+1])
					}
					clearIdxSet[idx] = struct{}{}
					clearIdxSet[idx+1] = struct{}{}
				}
			}
			if strings.HasPrefix(arg, key) {
				clearIdxSet[idx] = struct{}{}
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

	if g.Debug {
		g.Logger.Printf("To be cleaned %v", clear)
	}
	initialArgsSize := len(args)
	for _, idx := range clear {
		if g.Debug {
			g.Logger.Printf("Clearing %d", idx)
		}
		offset := initialArgsSize - len(args)
		offsetIdx := idx - offset
		if g.Debug {
			g.Logger.Printf("Actually clearing %d", offsetIdx)
		}
		if initialArgsSize-1 == idx {
			if g.Debug {
				g.Logger.Printf("Cut last arg")
			}
			args = args[:offsetIdx]
		} else {
			args = append(args[:offsetIdx], args[offsetIdx+1:]...)
		}
		if g.Debug {
			g.Logger.Printf("new slice is %v", args)
		}
	}
	g.Args = args
}

func (g *GccBuilder) parseCommandLine(args []string) {
	if g.Debug {
		g.Logger.Printf("Parsing arguments: %v", args)
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

	if g.Debug {
		g.Logger.Printf("Parsing cleaned commandline: %v", g.Args)
	}
	err := gccFlags.Parse(g.Args)
	if err != nil {
		g.Logger.Fatalf("Unrecoverable commandline parsing error: %s", err)
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
	if g.Debug {
		g.Logger.Printf("Mode set to: %v", g.Mode)
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

const (
	libPathVar = "LIBRARY_PATH"
)

// FindActualLibraries discovers the actual libraries on the path
func FindActualLibraries(libs []string, libpath []string) ([]string, error) {
	actualLibPaths := []string{}
	libpathvar, present := os.LookupEnv(libPathVar)
	if present && libpathvar != "" {
		libpath = append([]string{libpathvar}, libpath...)
	}
	var libprefix string
	var libsuffix []string
	var syslibpath []string
	switch runtime.GOOS {
	case "linux":
		libprefix = "lib"
		libsuffix = []string{".so", ".a"}
		syslibpath = []string{"/lib", "/usr/lib", "/usr/local/lib", "/lib64"}
	case "darwin":
		libprefix = "lib"
		libsuffix = []string{".dylib", ".so", ".a"}
		syslibpath = []string{"/usr/lib", "/usr/local/lib"}
	case "windows":
		libprefix = ""
		libsuffix = []string{".dll"}
		syslibpath = []string{""}
	}
	for _, lib := range libs {
		for _, dir := range append(libpath, syslibpath...) {
			if dir == "" {
				// Unix shell semantics: path element "" means "."
				dir = "."
			}
			err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
				for _, suffix := range libsuffix {
					if f.Name() == fmt.Sprintf("%s%s%s", libprefix, lib, suffix) {
						actualLibPaths = append(actualLibPaths, path)
						return fmt.Errorf("Found %s", path)
					}
				}
				return nil

			})
			if err != nil {
				break
			}
		}
	}

	if len(actualLibPaths) == len(libs) {
		return actualLibPaths, nil
	}

	return actualLibPaths, fmt.Errorf("Missing libraries from %v in %v", libs, actualLibPaths)
}
