package metrics

import "github.com/rai-project/dlframework"

/*
Compute the jaccard overlap of two sets of boxes.  The jaccard overlap
is simply the intersection over union of two boxes.  Here we operate on
ground truth boxes and default boxes.
E.g.:
    A ∩ B / A ∪ B = A ∩ B / (area(A) + area(B) - A ∩ B)
Args:
    box_a: Predicted bounding boxes
    box_b: Ground Truth bounding boxes
Return:
    jaccard overlap: Shape: [n_pred, n_gt]
*/

func BoundingBoxJaccard(boxA, boxB *dlframework.BoundingBox) float64 {
	intersection := BoundingBoxIntersectionOverUnion(boxA, boxA)
	areaA := float64(boxA.Area())
	areaB := float64(boxB.Area())
	union := areaA + areaB - intersection
	return intersection / union
}

func Jaccard(featA, featB *dlframework.Feature) float64 {
	boxA, ok := featA.Feature.(*dlframework.Feature_BoundingBox)
	if !ok {
		panic("unable to convert first feature to boundingbox")
	}
	boxB, ok := featB.Feature.(*dlframework.Feature_BoundingBox)
	if !ok {
		panic("unable to convert second feature to boundingbox")
	}
	return BoundingBoxJaccard(boxA.BoundingBox, boxB.BoundingBox)
}
