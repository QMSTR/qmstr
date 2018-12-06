package gnubuilder

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

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/spf13/afero"
)

const (
	libPathVar = "LIBRARY_PATH"
)

type Mode int

const (
	ModeLink Mode = iota
	ModePreproc
	ModeCompile
	ModeAssemble
	ModePrintOnly
	ModeUndef
)

func CleanCmdLine(args []string, logger *log.Logger, debug bool, staticLink bool, staticLibs map[string]struct{}, mode Mode) []string {
	clearIdxSet := map[int]struct{}{}
	for idx, arg := range args {

		if debug {
			logger.Printf("%d - %s", idx, arg)
		}

		// parse flags depending on the mode
		stringArgs := map[string]struct{}{}
		booleanArgs := map[string]struct{}{}
		switch mode {
		case ModeLink:
			stringArgs = LinkStringArgs
			booleanArgs = LinkBoolArgs
		case ModeAssemble:
			stringArgs = AssembleStringArgs
			booleanArgs = AssembleBoolArgs
		default:
			stringArgs = StringArgs
			booleanArgs = BoolArgs
		}
		// index string flags
		for key := range stringArgs {
			if idx < len(args)-1 {
				if debug {
					logger.Printf("Find %s string arg in %s with %s", key, fmt.Sprintf("%s %s ", arg, args[idx+1]), fmt.Sprintf("%s%s", key, StringArgsRE))
				}
				re := regexp.MustCompile(fmt.Sprintf("%s%s", key, StringArgsRE))
				if re.MatchString(fmt.Sprintf("%s %s ", arg, args[idx+1])) {
					if debug {
						logger.Printf("Found %v string arg", args[idx:idx+1])
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
		for key := range booleanArgs {
			if strings.HasPrefix(arg, key) {
				clearIdxSet[idx] = struct{}{}
				// use static libraries when linking statically and incrementally
				if arg == "-static" || arg == "-r" || arg == "--no-dynamic-linker" {
					staticLink = true
				}
				staticLib := StaticLibPattern.FindAllStringSubmatch(arg, 1)
				if staticLib != nil {
					staticLibs[staticLib[0][1]] = struct{}{}
				}

			}
		}

		// fix long arguments to pass through pflags
		for key := range FixPosixArgs {
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

	if debug {
		logger.Printf("To be cleaned %v", clear)
	}
	initialArgsSize := len(args)
	for _, idx := range clear {
		if debug {
			logger.Printf("Clearing %d", idx)
		}
		offset := initialArgsSize - len(args)
		offsetIdx := idx - offset
		if debug {
			logger.Printf("Actually clearing %d", offsetIdx)
		}
		if initialArgsSize-1 == idx {
			if debug {
				logger.Printf("Cut last arg")
			}
			args = args[:offsetIdx]
		} else {
			args = append(args[:offsetIdx], args[offsetIdx+1:]...)
		}
		if debug {
			logger.Printf("new slice is %v", args)
		}
	}
	return args
}

func CheckInputFileExt(inputFile string) service.FileNode_Type {
	ext := filepath.Ext(inputFile)
	switch ext {
	case ".c", ".cc", ".cpp", ".c++", ".cp", ".cxx":
		return service.FileNode_SOURCE
	case ".s", ".o", ".i", ".ii":
		return service.FileNode_INTERMEDIATE
	default:
		return service.FileNode_TARGET
	}
}

func GetOsLibFixes() (prefix string, dSuffixes []string, sSuffixes []string, err error) {
	switch runtime.GOOS {
	case "linux":
		return "lib", []string{".so"}, []string{".a"}, nil
	case "darwin":
		return "lib", []string{".dylib", ".so"}, []string{".a"}, nil
	case "windows":
		return "", []string{".dll"}, []string{".lib"}, nil
	}
	return "", nil, nil, errors.New("Platform not supported")
}

func GetSysLibPath() []string {
	switch runtime.GOOS {
	case "linux":
		return []string{"/lib", "/usr/lib", "/usr/local/lib", "/lib64"}
	case "darwin":
		return []string{"/usr/lib", "/usr/local/lib"}
	}
	return []string{""}
}

// FindActualLibraries discovers the actual libraries on the path
func FindActualLibraries(afs afero.Fs, actualLibs map[string]string, linkLibs []string, libPath []string, staticLink bool, staticLibs map[string]struct{}) error {
	aferofs := &afero.Afero{Fs: afs}
	libpathvar, present := os.LookupEnv(libPathVar)
	if present && libpathvar != "" {
		libPath = append([]string{libpathvar}, libPath...)
	}
	libprefix, libsuffix, staticSuffixes, err := GetOsLibFixes()
	if err != nil {
		return fmt.Errorf("could not search for libraries: %v", err)
	}

	// eliminate duplicated libs
	linkLibsSet := map[string]struct{}{}
	for _, lib := range linkLibs {
		linkLibsSet[lib] = struct{}{}
	}

	for _, dir := range libPath {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		aferofs.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			for lib := range linkLibsSet {
				// is lib located
				if _, ok := actualLibs[lib]; ok {
					continue
				}
				var suffixes []string
				// forced static lib or linking statically
				if _, ok := staticLibs[lib]; ok || staticLink {
					suffixes = staticSuffixes
				} else {
					suffixes = append(libsuffix, staticSuffixes...)
				}

				for _, suffix := range suffixes {
					if f.Name() == fmt.Sprintf("%s%s%s", libprefix, lib, suffix) {
						actualLibs[lib] = path
					}
				}
			}
			return nil
		})
	}

	if len(actualLibs) == len(linkLibsSet) {
		return nil
	}

	return fmt.Errorf("Missing libraries from %v in %v", linkLibs, actualLibs)
}
