package metrics

import (
	"testing"

	"github.com/rai-project/dlframework/framework/feature"
	"github.com/stretchr/testify/assert"
)

func TestIOU(t *testing.T) {
	boxA := feature.New(
		feature.BoundingBoxType(),
		feature.BoundingBoxXmin(39),
		feature.BoundingBoxXmax(203),
		feature.BoundingBoxYmin(63),
		feature.BoundingBoxYmax(112),
	)
	boxB := feature.New(
		feature.BoundingBoxType(),
		feature.BoundingBoxXmin(54),
		feature.BoundingBoxXmax(198),
		feature.BoundingBoxYmin(66),
		feature.BoundingBoxYmax(114),
	)

	iou := IntersectionOverUnion(boxA, boxB)

	assert.Equal(t, iou, 0.7957712638154734)
}
