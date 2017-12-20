package cmd

import (
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	framework "github.com/rai-project/dlframework/framework/cmd"
)

func forallmodels(run func() error) error {

	if modelName != "all" {
		return run()
	}

	outputDirectory := outputFileName
	if !com.IsDir(outputDirectory) {
		os.MkdirAll(outputDirectory, os.ModePerm)
	}
	for _, model := range framework.DefaultEvaulationModels {
		modelName, modelVersion = framework.ParseModelName(model)
		outputFileName = filepath.Join(outputDirectory, model+"."+outputFileExtension)
		run()
	}
	return nil
}
