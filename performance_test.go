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

func TestInsertPerformanceCollection(t *testing.T) {

	db, err := mongodb.NewDatabase(config.App.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()

	tbl, err := NewPerformanceCollection(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, tbl)

	defer tbl.Close()

	id := bson.NewObjectId()
	err = tbl.Insert(Performance{
		ID:        id,
		CreatedAt: time.Now(),
		// Framework: ....,
		// Model: ...,
		// Trace: ...,
		TraceLevel: tracer.CPU_ONLY_TRACE,
	})
	assert.NoError(t, err)
}

func TestInsertPerformanceCollectionFindByID(t *testing.T) {

	db, err := mongodb.NewDatabase(config.App.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()

	tbl, err := NewPerformanceCollection(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, tbl)

	defer tbl.Close()

	id := bson.NewObjectId()
	err = tbl.Insert(Performance{
		ID:         id,
		TraceLevel: tracer.CPU_ONLY_TRACE,
	})
	assert.NoError(t, err)

	var res []Performance
	err = tbl.Find(Performance{ID: id}, 0, 1, &res)
	assert.NoError(t, err)

	assert.Equal(t, res[0].TraceLevel, tracer.CPU_ONLY_TRACE)
}

func TestInsertPerformanceCollectionFindByTraceLevel(t *testing.T) {

	db, err := mongodb.NewDatabase(config.App.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()

	tbl, err := NewPerformanceCollection(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, tbl)

	defer tbl.Close()

	var res []Performance
	err = tbl.FindAll(Performance{TraceLevel: tracer.CPU_ONLY_TRACE}, &res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	for _, p := range res {
		assert.Equal(t, p.TraceLevel, tracer.CPU_ONLY_TRACE)
	}
}
