package main

type SamuraiDB struct {
	fileAdapter  *FileAdapter
	indexManager *IndexManager
}

func NewSamuraiDB(fileAdapter *FileAdapter, indexManager *IndexManager) *SamuraiDB {
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

//func main() {
//	dir := "./data"
//	fileAdapter := NewFileAdapter(dir)
//	indexManager := NewIndexManager(fileAdapter)
//	db := NewSamuraiDB(fileAdapter, indexManager)
//
//	err := db.Init()
//	if err != nil {
//		panic(err)
//	}
//
//	err = db.Set("key1", map[string]string{"data": "value"})
//	if err != nil {
//		panic(err)
//	}
//
//	value, err := db.Get("key1")
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("Retrieved value:", value)
//}
