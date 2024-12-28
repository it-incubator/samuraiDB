package indexmanager

import fa "samurai-db/internal/file-adapter"

type IndexManager struct {
	index       map[string]int64
	fileAdapter *fa.FileAdapter
}

func NewIndexManager(fa *fa.FileAdapter) *IndexManager {
	return &IndexManager{fileAdapter: fa}
}

func (im *IndexManager) Init() error {
	index, err := im.fileAdapter.ReadIndex()
	if err != nil {
		return err
	}
	im.index = make(map[string]int64)
	for key, value := range index {
		im.index[key] = value
	}
	return nil
}

func (im *IndexManager) SetOffset(key string, offset int64) error {
	im.index[key] = offset
	return im.fileAdapter.SaveIndex(im.index)
}

func (im *IndexManager) GetOffset(key string) (int64, bool) {
	offset, exists := im.index[key]
	return offset, exists
}
