package eventflow

const (
	// from https://github.com/williaster/data-ui/blob/master/packages/event-flow/src/constants.js
	// event attributes
	TS0     = "TS0"
	TS_NEXT = "TS_NEXT"
	TS_PREV = "TS_PREV"
	TS      = "TS"

	EVENT_NAME = "EVENT_NAME"
	ENTITY_ID  = "ENTITY_ID"
	EVENT_UUID = "EVENT_UUID"

	META            = "META"
	ELAPSED_MS_ROOT = "ELAPSED_MS_ROOT"
	ELAPSED_MS      = "ELAPSED_MS"
	EVENT_COUNT     = "EVENT_COUNT"

	ANY_EVENT_TYPE = "ANY_EVENT_TYPE"

	// scales
	ELAPSED_TIME_SCALE   = "ELAPSED_TIME_SCALE"
	EVENT_SEQUENCE_SCALE = "EVENT_SEQUENCE_SCALE"
	NODE_COLOR_SCALE     = "NODE_COLOR_SCALE"
	EVENT_COUNT_SCALE    = "EVENT_COUNT_SCALE"
	NODE_SEQUENCE_SCALE  = "NODE_SEQUENCE_SCALE"

	// node sorters
	ORDER_BY_EVENT_COUNT = "ORDER_BY_EVENT_COUNT"
	ORDER_BY_ELAPSED_MS  = "ORDER_BY_ELAPSED_MS"

	// note this can"t have spaces or it breaks its use in the Pattern "url(#id)""
	FILTERED_EVENTS = "FILTERED_EVENTS"
	CLIP_ID         = "CLIP_PATH_ID"
)
