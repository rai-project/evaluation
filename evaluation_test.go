package evaluation

import (
	"os"
	"testing"

	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
	"github.com/stretchr/testify/assert"

	"github.com/rai-project/config"
	"gopkg.in/mgo.v2"
)

func TestFindEvaluationsByModel(t *testing.T) {

	db, err := mongodb.NewDatabase("carml_step_trace")
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()

	tbl, err := NewEvaluationCollection(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, tbl)

	defer tbl.Close()

	evals, err := tbl.FindByModel(dlframework.ModelManifest{
		Name:    "BVLC-AlexNet",
		Version: "1.0",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, evals)

	for _, eval := range evals {
		assert.NotEmpty(t, eval.Model)
		assert.Equal(t, "BVLC-AlexNet", eval.Model.Name)
		assert.Equal(t, "1.0", eval.Model.Version)
	}
}

func TestMain(m *testing.M) {
	config.Init(
		config.AppName("carml"),
		config.VerboseMode(true),
		config.DebugMode(true),
	)
	mgo.SetDebug(true)
	mgo.SetLogger(log)
	os.Exit(m.Run())
}
