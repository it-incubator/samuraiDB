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

	return &FileAdapter{filename: filename, indexFileName: indexFileName}
}

func (fa *FileAdapter) Set(key string, data interface{}) (int64, error) {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	entry := fmt.Sprintf("%s,%s\n", key, serializeData(data))

	file, err := os.OpenFile(fa.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (fa *FileAdapter) Get(offset int64) (map[string]any, error) {
	if offset < 0 {
		return nil, fmt.Errorf("offset must be passed")
	}

	file, err := os.Open(fa.filename)
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

func (fa *FileAdapter) SaveIndex(indexMap map[string]int64) error {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	serializedMap, err := json.Marshal(indexMap)
	if err != nil {
		log.Fatal("were unable to serialize the data in SaveIndex")
	}
	return os.WriteFile(fa.indexFileName, serializedMap, 0644)
}

func (fa *FileAdapter) ReadIndex() (map[string]int64, error) {
	fileContent, err := os.ReadFile(fa.indexFileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return make(map[string]int64), nil
		}
		return nil, err
	}

	var index map[string]int64
	err = json.Unmarshal(fileContent, &index)
	if err != nil {
		return nil, err
	}

	return index, nil
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
