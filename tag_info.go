//+build ignore
//go:generate go get github.com/cheekybits/genny
//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "ElementType=uint,uint8,uint16,uint32,uint64,int,int8,int16,int32,int64,float32,float64,string"

package evaluation

import (
	"errors"

	"github.com/cheekybits/genny/generic"
	"github.com/spf13/cast"
	json "github.com/uber/jaeger/model/json"
)

type ElementType generic.Type

func getTagValueAsElementType(span json.Span, key string) (ElementType, error) {
	var res ElementType
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToElementTypeE(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsElementType(span json.Span, key string) ElementType {
	val, err := getTagValueAsElementType(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}
