package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileAdapter struct {
	filename      string
	indexFileName string
	mutex         sync.Mutex
}

func NewFileAdapter(dir string) *FileAdapter {
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

	offset, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	_, err = file.WriteString(entry)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func (fa *FileAdapter) Get(offset int64) (interface{}, error) {
	if offset < 0 {
		return nil, fmt.Errorf("Offset must be passed")
	}

	file, err := os.Open(fa.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Seek(offset, os.SEEK_SET)
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

	var value interface{}
	err = json.Unmarshal([]byte(storedValue), &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (fa *FileAdapter) SaveIndex(indexMap map[string]int64) error {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	serializedMap := serializeData(indexMap)
	return ioutil.WriteFile(fa.indexFileName, []byte(serializedMap), 0644)
}

func (fa *FileAdapter) ReadIndex() (map[string]interface{}, error) {
	fileContent, err := ioutil.ReadFile(fa.indexFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var index map[string]interface{}
	err = json.Unmarshal(fileContent, &index)
	if err != nil {
		return nil, err
	}

	return index, nil
}

// Helper functions
func serializeData(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func parseEntry(line string) (string, string, error) {
	i := 0
	for i < len(line) && line[i] != ',' {
		i++
	}
	if i >= len(line) {
		return "", "", fmt.Errorf("Invalid entry format")
	}
	return line[:i], line[i+1:], nil
}
