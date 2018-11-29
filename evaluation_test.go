package evaluation

import (
	"os"
	"testing"
	"time"

	"github.com/rai-project/config"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/tracer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		TraceLevel: tracer.FULL_TRACE,
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
		TraceLevel: tracer.FULL_TRACE,
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
		TraceLevel: tracer.FULL_TRACE,
	})
	assert.NoError(t, err)

	var res []Performance
	res, err = tbl.Find(Performance{ID: id}, 0, 1)
	assert.NoError(t, err)

	assert.Equal(t, res[0].TraceLevel, tracer.FULL_TRACE)
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
	err = tbl.FindAll(Performance{TraceLevel: tracer.FULL_TRACE}, &res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	for _, p := range res {
		assert.Equal(t, p.TraceLevel, tracer.FULL_TRACE)
	}
}

func TestFindEvaluationsByModel(t *testing.T) {
	db, err := mongodb.NewDatabase("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	defer db.Close()

	tbl, err := NewEvaluationCollection(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, tbl)

	defer tbl.Close()

	UserID := "test_UserID"
	evaluationEntry := Evaluation{
		UserID: UserID,
	}

	err = tbl.Insert(evaluationEntry)
	assert.NoError(t, err)

	evals, err := tbl.FindByUserID(UserID)
	assert.NoError(t, err)
	assert.NotEmpty(t, evals)

	for _, eval := range evals {
		assert.NotEmpty(t, eval.Model)
		assert.Equal(t, UserID, eval.UserID)
	}

	tbl.Delete()
	evals, err = tbl.FindByUserID(UserID)
	assert.Empty(t, evals)
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
