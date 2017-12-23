package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/GeertJohan/go-sourcepath"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"

	"github.com/rai-project/config"
	"github.com/rai-project/database"
	framework "github.com/rai-project/dlframework/framework/cmd"
	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"

	mongodb "github.com/rai-project/database/mongodb"
	_ "github.com/rai-project/logger/hooks"
	_ "github.com/rai-project/tracer/all"
)

var (
	limit                     int
	batchSize                 int
	goPath                    string
	mlArcWebAssetsPath        string
	raiSrcPath                string
	outputFileExtension       string
	hostName                  string
	machineArchitecture       string
	modelName                 string
	modelVersion              string
	frameworkName             string
	frameworkVersion          string
	databaseAddress           string
	databaseName              string
	databaseEndpoints         []string
	outputFileName            string
	outputFormat              string
	overwrite                 bool
	noHeader                  bool
	appendOutput              bool
	db                        database.Database
	evaluationCollection      *evaluation.EvaluationCollection
	performanceCollection     *evaluation.PerformanceCollection
	inputPerdictionCollection *evaluation.InputPredictionCollection
	modelAccuracyCollection   *evaluation.ModelAccuracyCollection
	divergenceCollection      *evaluation.DivergenceCollection

	sourcePath = sourcepath.MustAbsoluteDir()
)

func rootSetup() error {
	if databaseName == "" {
		databaseName = config.App.Name
	}
	if databaseAddress != "" {
		databaseEndpoints = []string{databaseAddress}
	}

	opts := []database.Option{}
	if len(databaseEndpoints) != 0 {
		opts = append(opts, database.Endpoints(databaseEndpoints))
	}

	var err error

	db, err = mongodb.NewDatabase(databaseName, opts...)
	if err != nil {
		return err
	}

	evaluationCollection, err = evaluation.NewEvaluationCollection(db)
	if err != nil {
		return err
	}

	performanceCollection, err = evaluation.NewPerformanceCollection(db)
	if err != nil {
		return err
	}

	modelAccuracyCollection, err = evaluation.NewModelAccuracyCollection(db)
	if err != nil {
		return err
	}

	inputPerdictionCollection, err = evaluation.NewInputPredictionCollection(db)
	if err != nil {
		return err
	}

	divergenceCollection, err = evaluation.NewDivergenceCollection(db)
	if err != nil {
		return err
	}

	inputPerdictionCollection, err = evaluation.NewInputPredictionCollection(db)
	if err != nil {
		return err
	}

	if outputFormat == "" && outputFileName != "" {
		outputFormat = filepath.Ext(outputFileName)
	}

	if fm, ok := framework.FrameworkNames[strings.ToLower(frameworkName)]; ok {
		frameworkName = fm
	}

	if modelName != "all" {
		outputFileExtension = filepath.Ext(outputFileName)
	} else {
		outputFileExtension = outputFormat
	}

	return nil
}

var EvaluationCmd = &cobra.Command{
	Use:   "evaluation",
	Short: "Get evaluation information from CarML",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running " + cmd.Name())
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		safeClose := func(cls ...io.Closer) {
			for _, c := range cls {
				if c == nil {
					return
				}
				c.Close()
			}
		}
		safeClose(
			evaluationCollection,
			performanceCollection,
			inputPerdictionCollection,
			modelAccuracyCollection,
			divergenceCollection,
			db,
		)

		return nil
	},
}

func init() {
	EvaluationCmd.PersistentFlags().StringVar(&hostName, "hostname", "", "hostname of machine to filter")
	EvaluationCmd.PersistentFlags().StringVar(&machineArchitecture, "arch", "", "architecture of machine to filter")
	EvaluationCmd.PersistentFlags().IntVar(&batchSize, "batch_size", 0, "the batch size to filter")

	EvaluationCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "modelName")
	EvaluationCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "modelVersion")
	EvaluationCmd.PersistentFlags().StringVar(&frameworkName, "framework_name", "", "frameworkName")
	EvaluationCmd.PersistentFlags().StringVar(&frameworkVersion, "framework_version", "", "frameworkVersion")
	EvaluationCmd.PersistentFlags().StringVar(&databaseAddress, "database_address", "", "address of the database")
	EvaluationCmd.PersistentFlags().StringVar(&databaseName, "database_name", "", "name of the database to query")

	EvaluationCmd.PersistentFlags().IntVar(&limit, "limit", -1, "limit the evaluations")
	EvaluationCmd.PersistentFlags().BoolVar(&overwrite, "overwrite", false, "if the file or directory exists, then they get deleted")
	EvaluationCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", "", "output file name")
	EvaluationCmd.PersistentFlags().BoolVar(&noHeader, "no_header", false, "show header labels for output")
	EvaluationCmd.PersistentFlags().BoolVar(&appendOutput, "append", false, "append the output")
	EvaluationCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "print format to use")

	EvaluationCmd.AddCommand(AllCmds...)
	EvaluationCmd.AddCommand(allCmd)

	pp.WithLineInfo = true
}

func init() {
	goPath = com.GetGOPATHs()[0]
	raiSrcPath = getSrcPath("github.com/rai-project")
	mlArcWebAssetsPath = filepath.Join(raiSrcPath, "ml-arc-web", "src", "assets", "data")
}
