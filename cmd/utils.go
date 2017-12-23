package cmd

import (
	"errors"
	"go/build"
	"path/filepath"

	"github.com/Unknwon/com"
)

func uptoIndex(arry []interface{}, idx int) int {
	if len(arry) <= idx {
		return len(arry) - 1
	}
	return idx
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func getSrcPath(importPath string) (appPath string) {
	paths := com.GetGOPATHs()
	for _, p := range paths {
		d := filepath.Join(p, "src", importPath)
		if com.IsExist(d) {
			appPath = d
			break
		}
	}

	if len(appPath) == 0 {
		appPath = filepath.Join(goPath, "src", importPath)
	}

	return appPath
}

func isExists(s string) bool {
	return com.IsExist(s)
}

func getBuildFile() (string, error) {
	pkg, err := build.Default.ImportDir(sourcePath, build.ImportMode(0))
	if err == nil && pkg.IsCommand() {
		return pkg, nil
	}

	mainPath := filepath.Join(sourcePath, "main.go")
	if com.IsFile(mainPath) {
		return mainPath, nil
	}

	return "", errors.New("unable to figure out what file to build")
}
