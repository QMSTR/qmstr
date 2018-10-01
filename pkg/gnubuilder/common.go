package gnubuilder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const (
	libPathVar = "LIBRARY_PATH"
)

func CleanCmdLine(args []string, logger *log.Logger, debug bool, staticLink bool, staticLibs map[string]struct{}, mode string) []string {
	clearIdxSet := map[int]struct{}{}
	for idx, arg := range args {

		if debug {
			logger.Printf("%d - %s", idx, arg)
		}

		// parse flags depending on the mode
		stringArgs := map[string]struct{}{}
		booleanArgs := map[string]struct{}{}
		switch mode {
		case "Link":
			stringArgs = LinkStringArgs
			booleanArgs = LinkBoolArgs
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

// FindActualLibraries discovers the actual libraries on the path
func FindActualLibraries(actualLibs map[string]string, linkLibs []string, libPath []string, staticLink bool, staticLibs map[string]struct{}) error {
	libpathvar, present := os.LookupEnv(libPathVar)
	if present && libpathvar != "" {
		libPath = append([]string{libpathvar}, libPath...)
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

	for _, dir := range append(libPath, syslibpath...) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			for _, lib := range linkLibs {
				// is lib located
				if _, ok := actualLibs[lib]; ok {
					continue
				}
				var suffixes []string
				// forced static lib or linking statically
				if _, ok := staticLibs[lib]; ok || staticLink {
					suffixes = []string{".a"}
				} else {
					suffixes = libsuffix
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

	if len(actualLibs) == len(linkLibs) {
		return nil
	}

	return fmt.Errorf("Missing libraries from %v in %v", linkLibs, actualLibs)
}
