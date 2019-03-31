package metrics

import (
	"github.com/chewxy/math32"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
)

// https://stackoverflow.com/questions/28723670/intersection-over-union-between-two-detections
// https://resources.wolframcloud.com/NeuralNetRepository/resources/SSD-VGG-300-Trained-on-PASCAL-VOC-Data
func BoundingBoxIntersectionOverUnion(boxA, boxB *dlframework.BoundingBox) float64 {

	// determine the (x, y)-coordinates of the intersection rectangle
	xA := math32.Max(boxA.GetXmin(), boxB.GetXmin())
	yA := math32.Max(boxA.GetYmin(), boxB.GetYmin())
	xB := math32.Min(boxA.GetXmax(), boxB.GetXmax())
	yB := math32.Min(boxA.GetYmax(), boxB.GetYmax())

	// compute the area of intersection rectangle
	interArea := float64(xB-xA) * float64(yB-yA)

	// compute the area of both the prediction and ground-truth
	// rectangles
	boxAArea := float64(boxA.GetXmax()-boxA.GetXmin()) * float64(boxA.GetYmax()-boxA.GetYmin())
	boxBArea := float64(boxB.GetXmax()-boxB.GetXmin()) * float64(boxB.GetYmax()-boxB.GetYmin())

	// compute the intersection over union by taking the intersection
	// area and dividing it by the sum of prediction + ground-truth
	// areas - the interesection area
	iou := interArea / (boxAArea + boxBArea - interArea)

	// return the intersection over union value
	return iou
}

func IntersectionOverUnion(featA, featB *dlframework.Feature) float64 {
	boxA, ok := featA.Feature.(*dlframework.Feature_BoundingBox)
	if !ok {
		panic("unable to convert first feature to boundingbox")
	}
	boxB, ok := featB.Feature.(*dlframework.Feature_BoundingBox)
	if !ok {
		panic("unable to convert second feature to boundingbox")
	}
	return BoundingBoxIntersectionOverUnion(boxA.BoundingBox, boxB.BoundingBox)
}

func init() {
	config.AfterInit(func() {
		RegisterFeatureCompareFunction("IntersectionOverUnion", IntersectionOverUnion)
	})
}
