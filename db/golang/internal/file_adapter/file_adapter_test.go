package fileadapter

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestNewAdapter(t *testing.T) {
	dir := t.TempDir()
	adapter := NewAdapter(dir)

	assert.Equal(t, filepath.Join(dir, "samuraidb.txt"), adapter.filename)
	assert.Equal(t, filepath.Join(dir, "index.txt"), adapter.indexFileName)
}

func TestFileAdapter_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	adapter := NewAdapter(dir)

	key := "exampleKey"
	value := map[string]any{"name": "Samurai"}
	offset, err := adapter.Set(key, value, 0)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, offset, int64(0))

	readValue, err := adapter.Get(offset, 0)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.ObjectsAreEqualValues(value, readValue)
	_, err = adapter.Get(offset, 1)
	assert.Error(t, err)
}

func TestFileAdapter_SaveAndReadIndex(t *testing.T) {
	dir := t.TempDir()
	adapter := NewAdapter(dir)

	value, err := json.Marshal(map[string]int64{"name": 123, "name2": 456})
	require.NoError(t, err)

	err = adapter.SaveIndexRaw(value)
	require.NoError(t, err)

	readIndex, err := adapter.ReadRawIndex()
	require.NoError(t, err)

	assert.Equal(t, value, readIndex)

	value, err = json.Marshal(map[string]int64{})
	require.NoError(t, err)

	err = adapter.SaveIndexRaw(value)
	require.NoError(t, err)

	readIndex, err = adapter.ReadRawIndex()
	require.NoError(t, err)

	var readIndexRes map[string]int64

	err = json.Unmarshal(readIndex, &readIndexRes)
	require.NoError(t, err)

	assert.Equal(t, value, readIndex)
}

func TestFileAdapter_GetInvalidOffset(t *testing.T) {
	dir := t.TempDir()
	adapter := NewAdapter(dir)

	_, err := adapter.Get(-1, 0)
	assert.Error(t, err)
	assert.Equal(t, "offset must be passed", err.Error())
}

func TestFileAdapter_TestParseEntry(t *testing.T) {

	line := "key,{value: 1231, value1: 'bla-bla'}"
	key, value, err := parseEntry(line)
	assert.NoError(t, err)
	assert.Equal(t, "key", key)
	assert.Equal(t, "{value: 1231, value1: 'bla-bla'}", value)

	line = ""
	_, _, err = parseEntry(line)
	assert.Error(t, err)

	line = "asdfhjklljh"
	_, _, err = parseEntry(line)
	assert.Error(t, err)
}
