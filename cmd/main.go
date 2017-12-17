// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/rai-project/config"
	"github.com/rai-project/dlframework/framework/cmd"
	evalcmd "github.com/rai-project/evaluation/cmd"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/evaluate")
)

func main() {
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/evaluate")
	})
	cmd.Init()

	if err := evalcmd.EvaluationCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
