package cmd

import (
	"io"
	"path"

	"github.com/rai-project/config"
	"github.com/rai-project/database"
	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"

	mongodb "github.com/rai-project/database/mongodb"
	_ "github.com/rai-project/logger/hooks"
	_ "github.com/rai-project/tracer/all"
)

var (
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
	db                        database.Database
	evaluationCollection      *evaluation.EvaluationCollection
	performanceCollection     *evaluation.PerformanceCollection
	inputPerdictionCollection *evaluation.InputPredictionCollection
	modelAccuracyCollection   *evaluation.ModelAccuracyCollection
	divergenceCollection      *evaluation.DivergenceCollection
)

var EvaluationCmd = &cobra.Command{
	Use:   "evaluation",
	Short: "Get evaluation information from CarML",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
			outputFormat = path.Ext(outputFileName)
		}

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

	EvaluationCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "modelName")
	EvaluationCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "modelVersion")
	EvaluationCmd.PersistentFlags().StringVar(&frameworkName, "framework_name", "", "frameworkName")
	EvaluationCmd.PersistentFlags().StringVar(&frameworkVersion, "framework_version", "", "frameworkVersion")
	EvaluationCmd.PersistentFlags().StringVar(&databaseAddress, "database_address", "", "address of the database")
	EvaluationCmd.PersistentFlags().StringVar(&databaseName, "database_name", "", "name of the database to query")

	EvaluationCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", "", "output file name")
	EvaluationCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "print format to use")

	EvaluationCmd.AddCommand(durationCmd)
	EvaluationCmd.AddCommand(latencyCmd)
}
