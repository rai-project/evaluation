package evaluation

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

func uptoIndex(arry []interface{}, idx int) int {
	if len(arry) <= idx {
		return len(arry) - 1
	}
	return idx
}

func toFloat64Slice(i interface{}) []float64 {
	res, _ := toFloat64SliceE(i)
	return res
}

func toFloat64SliceE(i interface{}) ([]float64, error) {
	if i == nil {
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}

	switch v := i.(type) {
	case []float64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat64E(s.Index(j).Interface())
			if err != nil {
				return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}
}

func float64SliceToStringSlice(us []float64) []string {
	res := make([]string, len(us))
	for ii, u := range us {
		res[ii] = cast.ToString(u)
	}
	return res
}

func uint64SliceToStringSlice(us []uint64) []string {
	res := make([]string, len(us))
	for ii, u := range us {
		res[ii] = cast.ToString(u)
	}
	return res
}

func predictIndexOf(span model.Span, predictSpans Spans) int {
	for ii, predict := range predictSpans {
		if span.ParentSpanID == predict.SpanID {
			return ii
		}
		for _, ref := range span.References {
			if ref.RefType == model.ChildOf && ref.SpanID == predict.SpanID {
				return ii
			}
		}
	}
	return -1
}

func tagsOf(span model.Span) map[string]string {
	res := map[string]string{}
	for _, lg := range span.Logs {
		for _, fld := range lg.Fields {
			res[fld.Key] = cast.ToString(fld.Value)
		}
	}
	for _, tag := range span.Tags {
		res[tag.Key] = cast.ToString(tag.Value)
	}
	return res
}

func parentOf(span model.Span) model.SpanID {
	if span.ParentSpanID != "" {
		return span.ParentSpanID
	}
	for _, ref := range span.References {
		if ref.RefType == model.ChildOf {
			return ref.SpanID
		}
	}
	return model.SpanID("")
}
