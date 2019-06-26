package cmd

import (
	"errors"
	"go/build"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/rai-project/evaluation"
	udb "upper.io/db.v3"
)

func getEvaluations() (evaluation.Evaluations, error) {

	filter := udb.Cond{}
	if modelName != "" {
		filter["model.name"] = modelName
	}
	if modelVersion != "" {
		filter["model.version"] = modelVersion
	}
	if frameworkName != "" {
		filter["framework.name"] = frameworkName
	}
	if frameworkVersion != "" {
		filter["framework.version"] = frameworkVersion
	}
	if machineArchitecture != "" {
		filter["machinearchitecture"] = machineArchitecture
	}
	if hostName != "" {
		filter["hostname"] = hostName
	}
	if batchSize != 0 {
		filter["batch_size"] = batchSize
	}
	evals, err := evaluationCollection.Find(filter)
	if err != nil {
		return nil, err
	}

	if limit > 0 {
		evals = evals[:minInt(len(evals)-1, limit)]
	}

	return evaluation.Evaluations(evals), nil
}

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
		return pkg.SrcRoot, nil
	}

	mainPath := filepath.Join(sourcePath, "main.go")
	if com.IsFile(mainPath) {
		return mainPath, nil
	}

	return "", errors.New("unable to figure out what file to build")
}
