package indexmanager

import (
	"github.com/stretchr/testify/assert"
	"samurai-db/internal/file_adapter"
	"testing"
)

func TestIndexManager(t *testing.T) {

	tests := []struct {
		name     string
		testFunc func(*testing.T, *IndexManager)
	}{
		{
			name: "Test Init function",
			testFunc: func(t *testing.T, indexManager *IndexManager) {
				err := indexManager.Init()
				assert.NoError(t, err)

				expectedIndex := make(map[string]Index)
				assert.Equal(t, expectedIndex, indexManager.index)
			},
		},
		{
			name: "Test SetOffset function for segment 0",
			testFunc: func(t *testing.T, indexManager *IndexManager) {
				key := "test_key"
				offset := Index{0, 100}
				err := indexManager.SetOffset(key, offset.Offset, offset.Segment)
				assert.NoError(t, err)

				newOffset, exists := indexManager.GetOffset(key)
				assert.True(t, exists)
				assert.Equal(t, offset, newOffset)
			},
		},
		{
			name: "Test SetOffset function for segment 1",
			testFunc: func(t *testing.T, indexManager *IndexManager) {
				key := "test_key"
				offset := Index{1, 200}
				err := indexManager.SetOffset(key, offset.Offset, offset.Segment)
				assert.NoError(t, err)

				newOffset, exists := indexManager.GetOffset(key)
				assert.True(t, exists)
				assert.Equal(t, offset, newOffset)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			fileAdapter := fileadapter.NewAdapter(dir)

			indexManager := NewIndexManager(fileAdapter)
			err := indexManager.Init()
			assert.NoError(t, err)

			tc.testFunc(t, indexManager)
		})
	}
}
