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
	evals, err := evaluationCollection.Find(filter)
	if err != nil {
		return nil, err
	}

	return evaluation.Evaluations(evals), nil
}
