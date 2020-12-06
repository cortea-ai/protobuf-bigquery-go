package protobq

import (
	"testing"

	"cloud.google.com/go/bigquery"
	"google.golang.org/genproto/googleapis/example/library/v1"
	"gotest.tools/v3/assert"
)

func TestMessageSaver_Save(t *testing.T) {
	messageSaver := MessageSaver{
		Message: &library.Book{
			Name:   "test-name",
			Author: "test-author",
			Title:  "test-title",
			Read:   true,
		},
		InsertID: "test-insert-id",
	}
	actual, insertID, err := messageSaver.Save()
	assert.NilError(t, err)
	assert.Equal(t, messageSaver.InsertID, insertID)
	expected := map[string]bigquery.Value{
		"name":   "test-name",
		"author": "test-author",
		"title":  "test-title",
		"read":   true,
	}
	assert.DeepEqual(t, expected, actual)
}
