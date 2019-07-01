package evaluation

import (
	"github.com/pkg/errors"
	model "github.com/uber/jaeger/model/json"
)

func findPredictStep(spans []model.Span) (model.Span, error) {
	return findSpanByOperationName(spans, "PredictStep")
}

func findCPredict(spans []model.Span) (model.Span, error) {
	return findSpanByOperationName(spans, "c_predict")
}

func findModelName(spans []model.Span) (string, error) {
	predictSpan, err := findPredictStep(spans)
	if err != nil {
		return "", err
	}
	return getTagValueAsString(predictSpan, "model_name")
}

func findBatchSize(spans []model.Span) (int, error) {
	predictSpan, err := findPredictStep(spans)
	if err != nil {
		return 0, err
	}
	return getTagValueAsInt(predictSpan, "batch_size")
}

func findSpanByOperationName(spans []model.Span, operationName string) (model.Span, error) {
	for _, span := range spans {
		if span.OperationName == operationName {
			return span, nil
		}
	}
	return model.Span{}, errors.Errorf("the span with operationName = %v was not found", operationName)
}
