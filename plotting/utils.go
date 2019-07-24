package plotting

import (
	"errors"
	"path/filepath"
	"sort"

	"github.com/mattn/go-zglob"
	"github.com/spf13/cast"
)

func contains(lst interface{}, elem interface{}) bool {
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
	paths, err := zglob.Glob(path)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.StringSlice(paths))
	return paths, nil
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
	paths, err := zglob.Glob(path)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.StringSlice(paths))
	return paths, nil
}
