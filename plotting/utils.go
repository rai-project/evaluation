package plotting

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/mattn/go-zglob"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

func contains(lst []interface{}, elem interface{}) bool {
	return funk.Contains(lst, elem)
}

func fullTraceFile(o *Options) (string, error) {
	return evaluationTraceFile(o, "trace_full_trace.json")
}

func fullTraceFiles(o *Options) ([]string, error) {
	return evaluationTraceFiles(o, "trace_full_trace.json")
}

func modelTraceFile(o *Options) (string, error) {
	return evaluationTraceFile(o, "trace_model_trace.json")
}

func modelTraceFiles(o *Options) ([]string, error) {
	return evaluationTraceFiles(o, "trace_model_trace.json")
}

func frameworkTraceFile(o *Options) (string, error) {
	return evaluationTraceFile(o, "trace_framework_trace.json")
}

func frameworkTraceFiles(o *Options) ([]string, error) {
	return evaluationTraceFiles(o, "trace_framework_trace.json")
}

func evaluationTraceFile(o *Options, filename string) (string, error) {
	path, err := evaluationTraceDir(o)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, filename), nil
}

func evaluationTraceFiles(o *Options, filename string) ([]string, error) {
	paths, err := evaluationTraceDirs(o)
	if err != nil {
		return nil, err
	}
	res := make([]string, len(paths))
	for ii, path := range paths {
		res[ii] = filepath.Join(path, filename)
	}
	return res, nil
}

func evaluationTraceDir(o *Options) (string, error) {
	files, err := evaluationTraceDirs(o)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", errors.New("too mancy files")
	}
	return files[0], nil
}

func evaluationTraceDirs(o *Options) ([]string, error) {
	orGlob := func(s string) string {
		if s == "" {
			return "**"
		}
		return s
	}
	orGlobInt := func(s int) string {
		if s == 0 {
			return "**"
		}
		return cast.ToString(s)
	}
	orGlobBool := func(s *bool) string {
		if s == nil {
			return "**"
		}
		if *s == true {
			return "gpu"
		}
		return "cpu"
	}
	path := filepath.Join(
		orGlob(o.baseDir),
		orGlob(o.frameworkName),
		orGlob(o.frameworkVersion),
		orGlob(o.modelName),
		orGlob(o.modelVersion),
		orGlobInt(o.batchSize),
		orGlobBool(o.useGPU),
		orGlob(o.machineHostName),
	)
	return zglob.Glob(path)
}

func evaluationDirs(o *Options) ([]string, error) {
	orGlob := func(s string) string {
		if s == "" {
			return "**"
		}
		return s
	}
	path := filepath.Join(
		orGlob(o.baseDir),
		orGlob(o.frameworkName),
		orGlob(o.frameworkVersion),
		orGlob(o.modelName),
		orGlob(o.modelVersion),
	)
	return zglob.Glob(path)
}

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// The filename is generated by taking pattern and adding a random
// string to the end. If pattern includes a "*", the random string
// replaces the last "*".
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func tempFile(dir, pattern string) string {
	if dir == "" {
		dir = os.TempDir()
	}

	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}

	var name string

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name = filepath.Join(dir, prefix+nextRandom()+suffix)
		if com.IsFile(name) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	if !com.IsDir(filepath.Dir(name)) {
		os.MkdirAll(filepath.Dir(name), os.ModePerm)
	}
	return name
}
