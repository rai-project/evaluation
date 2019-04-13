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

func TestResnetBatches(t *testing.T) {
	batch, err := NewBatchPlot("ResNet GPU",
		Option.UseGPU(true),
		Option.ModelName("ResNet*"),
		Option.IgnoreReadErrors(true),
	)
	assert.NoError(t, err)
	assert.NotNil(t, batch)

	err = batch.Open()
	assert.NoError(t, err)

	// pp.Println(batch)
}

func TestResnetBatches(t *testing.T) {
	batch, err := NewBatchPlot("ResNet CPU",
		Option.UseGPU(false),
		Option.ModelName("ResNet*"),
		Option.IgnoreReadErrors(true),
	)
	assert.NoError(t, err)
	assert.NotNil(t, batch)

	err = batch.Open()
	assert.NoError(t, err)

	// pp.Println(batch)
}

func TestMobileNetBatches(t *testing.T) {
	batch, err := NewBatchPlot("Image", Option.ModelName("MobileNet*"), Option.IgnoreReadErrors(true))
	assert.NoError(t, err)
	assert.NotNil(t, batch)

	err = batch.Open()
	assert.NoError(t, err)

	// pp.Println(batch)
}
