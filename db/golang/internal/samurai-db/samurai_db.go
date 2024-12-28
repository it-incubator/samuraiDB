package samuraidb

import (
	fa "samurai-db/internal/file_adapter"
	im "samurai-db/internal/index_manager"
)

type SamuraiDB struct {
	fileAdapter  *fa.FileAdapter
	indexManager *im.IndexManager
}

func NewSamuraiDB(fileAdapter *fa.FileAdapter, indexManager *im.IndexManager) *SamuraiDB {
	return &SamuraiDB{fileAdapter: fileAdapter, indexManager: indexManager}
}

func (db *SamuraiDB) Init() error {
	return db.indexManager.Init()
}

func (db *SamuraiDB) Set(key string, data interface{}) error {
	offset, err := db.fileAdapter.Set(key, data)
	if err != nil {
		return err
	}
	return db.indexManager.SetOffset(key, offset)
}

func (db *SamuraiDB) Get(key string) (interface{}, error) {
	offset, exists := db.indexManager.GetOffset(key)
	if !exists {
		return nil, nil // Key not found
	}
	return db.fileAdapter.Get(offset)
}
