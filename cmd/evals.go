package cmd

import (
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
