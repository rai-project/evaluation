package plotting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	batch, err := NewBatchPlot("test", Option.IgnoreReadErrors(true))
	assert.NoError(t, err)
	assert.NotNil(t, batch)

	err = batch.Open()
	assert.NoError(t, err)

	// pp.Println(batch)
}
