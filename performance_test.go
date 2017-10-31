package evaluation

import (
	"testing"
	"time"

	"github.com/rai-project/tracer"
	"gopkg.in/mgo.v2/bson"

	"github.com/rai-project/config"
	mongodb "github.com/rai-project/database/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	db, err := mongodb.NewDatabase(config.App.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()
}

func TestInsertPerformance(t *testing.T) {
	db, err := mongodb.NewDatabase(config.App.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()

	tbl, err := mongodb.NewTable(db, "performance")
	assert.NoError(t, err)
	assert.NotEmpty(t, tbl)

	tbl.Create(nil)

	err = tbl.Insert(Performance{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		// Framework: ....,
		// Model: ...,
		// Trace: ...,
		TraceLevel: tracer.CPU_ONLY_TRACE,
	})
	assert.NoError(t, err)

}
