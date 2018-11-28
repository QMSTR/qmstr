package gnubuilder

import "regexp"

var (
	BoolArgs = map[string]struct{}{
		"-w":                struct{}{},
		"-W":                struct{}{},
		"-O":                struct{}{},
		"-f":                struct{}{},
		"-C":                struct{}{},
		"-std":              struct{}{},
		"-nostdinc":         struct{}{},
		"-nostdlib":         struct{}{},
		"-print-file-name":  struct{}{},
		"--print-file-name": struct{}{},
		"-M":                struct{}{},
		"-MG":               struct{}{},
		"-MM":               struct{}{},
		"-MD":               struct{}{},
		"-MMD":              struct{}{},
		"-MP":               struct{}{},
		"-v":                struct{}{},
		"-g":                struct{}{},
		"-pg":               struct{}{},
		"-P":                struct{}{},
		"-pipe":             struct{}{},
		"-pie":              struct{}{},
		"-pedantic":         struct{}{},
		"-print-":           struct{}{},
		"-pthread":          struct{}{},
		"-rdynamic":         struct{}{},
		"-shared":           struct{}{},
		"-static":           struct{}{},
		"-dynamiclib":       struct{}{},
		"-dumpversion":      struct{}{},
		"-dM":               struct{}{},
		"--version":         struct{}{},
		"-undef":            struct{}{},
		"-nostartfiles":     struct{}{},
		"-remap":            struct{}{},
		"-r":                struct{}{},
	}

	StringArgs = map[string]struct{}{
		"-e":                     struct{}{},
		"-D":                     struct{}{},
		"-B":                     struct{}{},
		"-Q":                     struct{}{},
		"-T":                     struct{}{},
		"-U":                     struct{}{},
		"-m":                     struct{}{},
		"-x":                     struct{}{},
		"-MF":                    struct{}{},
		"-MT":                    struct{}{},
		"-MQ":                    struct{}{},
		"-install_name":          struct{}{},
		"-compatibility_version": struct{}{},
		"-current_version":       struct{}{},
		"--param":                struct{}{},
	}

	StringArgsRE = "\\s+\\S+={0,1}\\S*\\s"

	FixPosixArgs = map[string]struct{}{
		"-isystem": struct{}{},
		"-include": struct{}{},
	}

	StaticLibPattern = regexp.MustCompile("-static-lib(\\w+)")

	LinkBoolArgs = map[string]struct{}{
		"-P":                  struct{}{},
		"-E":                  struct{}{},
		"-F":                  struct{}{},
		"-g":                  struct{}{},
		"-r":                  struct{}{},
		"-i":                  struct{}{},
		"-q":                  struct{}{},
		"-v":                  struct{}{},
		"-pie":                struct{}{},
		"-static":             struct{}{},
		"-shared":             struct{}{},
		"-nostdlib":           struct{}{},
		"--emit-relocs":       struct{}{},
		"--whole-archive":     struct{}{},
		"--no-whole-archive":  struct{}{},
		"--no-undefined":      struct{}{},
		"--start-group":       struct{}{},
		"--end-group":         struct{}{},
		"--no-dynamic-linker": struct{}{},
		"-dynamic-linker":     struct{}{},
		"--push-state":        struct{}{},
		"--pop-state":         struct{}{},
		"--build-id":          struct{}{},
		"--eh-frame-hdr":      struct{}{},
		"--help":              struct{}{},
		"--as-needed":         struct{}{},
		"--gc-sections":       struct{}{},
		"--no-gc-sections":    struct{}{},
		"--noexecstack":       struct{}{},
	}

	LinkStringArgs = map[string]struct{}{
		"-B":           struct{}{},
		"-T":           struct{}{},
		"-m":           struct{}{},
		"-z":           struct{}{},
		"-plugin":      struct{}{},
		"-h":           struct{}{},
		"-soname":      struct{}{},
		"--sysroot":    struct{}{},
		"--hash-style": struct{}{},
	}

	AssembleStringArgs = map[string]struct{}{}

	AssembleBoolArgs = map[string]struct{}{
		"--32":          struct{}{},
		"--64":          struct{}{},
		"--noexecstack": struct{}{},
		"-v":            struct{}{},
	}
)
