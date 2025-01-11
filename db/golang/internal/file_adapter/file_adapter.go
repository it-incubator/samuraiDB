package fileadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileAdapter struct {
	dir           string
	filename      string
	indexFileName string
	mutex         sync.Mutex
}

func NewAdapter(dir string) *FileAdapter {
	filename := filepath.Join(dir, "samuraidb.txt")
	indexFileName := filepath.Join(dir, "index.txt")

	// Проверка или создание директории
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	return &FileAdapter{dir: dir, filename: filename, indexFileName: indexFileName}
}

func (fa *FileAdapter) Set(key string, data any, segment int64) (int64, error) {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	entry := fa.StringifyEntry(key, data)

	filename, _ := fa.getSegmentFilename(segment)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = file.WriteString(entry)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func (fa *FileAdapter) Get(offset, segment int64) (map[string]any, error) {
	if offset < 0 {
		return nil, fmt.Errorf("offset must be passed")
	}

	filename, _ := fa.getSegmentFilename(segment)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	line := string(buffer[:n])
	_, storedValue, _ := parseEntry(line)

	var value map[string]any
	err = json.Unmarshal([]byte(storedValue), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (fa *FileAdapter) SaveIndexRaw(indexMapRaw []byte) error {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	return os.WriteFile(fa.indexFileName, indexMapRaw, 0644)
}

func (fa *FileAdapter) ReadRawIndex() ([]byte, error) {
	fileContent, err := os.ReadFile(fa.indexFileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []byte{}, nil
		}
		return nil, err
	}

	return fileContent, nil
}

func (fa *FileAdapter) GetFileSize(segment int64) (int64, error) {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	filename, _ := fa.getSegmentFilename(segment)

	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("не удалось получить информацию о файле: %v", err)
	}

	return fileInfo.Size(), nil
}

// Helper functions
func serializeData(data any) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func parseEntry(line string) (string, string, error) {
	res := strings.SplitN(line, ",", 2)
	if len(res) < 2 {
		return "", "", fmt.Errorf("invalid entry format")
	}

	return res[0], res[1], nil
}

func (fa *FileAdapter) StringifyEntry(key string, data any) string {
	entry := fmt.Sprintf("%s,%s\n", key, serializeData(data))

	return entry
}

func (fa *FileAdapter) getSegmentFilename(segment int64) (string, error) {
	segmentPath := filepath.Join(fa.dir, fmt.Sprintf("samuraidb_segment_%d.txt", segment))
	return segmentPath, nil
}
