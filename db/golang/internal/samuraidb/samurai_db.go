package samuraidb

import (
	im "samurai-db/internal/index_manager"
	sm "samurai-db/internal/segment_manager"
)

//type SegmentManager interface {
//	Set(key string, data any) (int64, int64, error)
//	Get(offset, segment int64) (map[string]any, error)
//}

type SamuraiDB struct {
	segmentManager *sm.SegmentManager // use interface SegmentManager
	indexManager   *im.IndexManager
}

func NewSamuraiDB(fileAdapter *sm.SegmentManager, indexManager *im.IndexManager) *SamuraiDB {
	return &SamuraiDB{segmentManager: fileAdapter, indexManager: indexManager}
}

func (db *SamuraiDB) Init() error {
	return db.indexManager.Init()
}

func (db *SamuraiDB) Set(key string, data interface{}) error {
	offset, segment, err := db.segmentManager.Set(key, data)
	if err != nil {
		return err
	}
	return db.indexManager.SetOffset(key, offset, segment)
}

func (db *SamuraiDB) Get(key string) (map[string]any, error) {
	indexData, exists := db.indexManager.GetOffset(key)
	if !exists {
		return nil, nil
	}

	get, err := db.segmentManager.Get(indexData.Offset, indexData.Segment)
	if err != nil {
		return nil, err
	}

	return get, nil
}
