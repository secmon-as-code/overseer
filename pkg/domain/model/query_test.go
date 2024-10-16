package model_test

import (
	"embed"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
)

//go:embed testdata/query/*.sql
var queryFiles embed.FS

func TestQuery(t *testing.T) {
	type testCase struct {
		fileName string
		isErr    bool
		validate func(t *testing.T, q *model.Query)
	}

	tf := func(tc testCase) func(t *testing.T) {
		return func(t *testing.T) {
			data, err := queryFiles.ReadFile(tc.fileName)
			gt.NoError(t, err)

			q, err := model.NewQuery(data)
			if tc.isErr {
				gt.Error(t, err)
				return
			}

			gt.NoError(t, err)
			tc.validate(t, q)
		}
	}

	t.Run("valid cases", func(t *testing.T) {
		t.Run("simple", tf(testCase{
			fileName: "testdata/query/valid1.sql",
			validate: func(t *testing.T, q *model.Query) {
				gt.Equal(t, "test1", q.ID())
			},
		}))

		t.Run("trim line", tf(testCase{
			fileName: "testdata/query/valid2.sql",
			validate: func(t *testing.T, q *model.Query) {
				gt.Equal(t, "test2", q.ID())
			},
		}))
	})

	t.Run("invalid cases", func(t *testing.T) {
		t.Run("metadata in comment out line", tf(testCase{
			fileName: "testdata/query/invalid1.sql",
			isErr:    true,
		}))

		t.Run("no metadata", tf(testCase{
			fileName: "testdata/query/invalid2.sql",
			isErr:    true,
		}))
	})
}
