package indexmanager

import (
	"encoding/json"
	"log"
	fa "samurai-db/internal/file_adapter"
)

type Index struct {
	Segment int64 `json:"segment"`
	Offset  int64 `json:"offset"`
}

type IndexManager struct {
	index       map[string]Index
	fileAdapter *fa.FileAdapter
}

func NewIndexManager(fa *fa.FileAdapter) *IndexManager {
	return &IndexManager{fileAdapter: fa}
}

func (im *IndexManager) Init() error {
	rawData, err := im.fileAdapter.ReadRawIndex()
	if err != nil {
		return err
	}

	if len(rawData) == 0 {
		im.index = make(map[string]Index)
		return nil
	}

	var oldFormatData map[string]int64
	if err = json.Unmarshal(rawData, &oldFormatData); err == nil {
		im.index = make(map[string]Index)
		for key, offset := range oldFormatData {
			im.index[key] = Index{
				Segment: 0,
				Offset:  offset,
			}
		}

		log.Printf("Используется старый формат индексации")
		// todo: Перезаписать данные на новый формат
		return nil
	}

	var indexData map[string]Index
	if err = json.Unmarshal(rawData, &indexData); err != nil {
		return err
	}

	im.index = indexData
	return nil
}

func (im *IndexManager) SetOffset(key string, offset, segment int64) error {
	im.index[key] = Index{
		segment,
		offset,
	}

	serializedMap, err := json.Marshal(im.index)
	if err != nil {
		log.Fatal("were unable to serialize the data in SaveIndexRaw")
	}

	return im.fileAdapter.SaveIndexRaw(serializedMap)
}

func (im *IndexManager) GetOffset(key string) (Index, bool) {
	offset, exists := im.index[key]
	return offset, exists
}
