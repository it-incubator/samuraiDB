package segmentmanager

import (
	"fmt"

	fa "samurai-db/internal/file_adapter"
)

type SegmentManager struct {
	currentSegment int64
	segmentSize    int64
	fileAdapter    *fa.FileAdapter
}

func NewSegmentManager(fa *fa.FileAdapter) *SegmentManager {
	return &SegmentManager{fileAdapter: fa, currentSegment: 0, segmentSize: 1024}
}

func (s *SegmentManager) Set(key string, data any) (int64, int64, error) {
	foundCurrentSegmentFolder := true

	// fixme: Костыльный поиск последнего сегмента
	for foundCurrentSegmentFolder {
		size, _ := s.fileAdapter.GetFileSize(s.currentSegment)

		entry := s.fileAdapter.StringifyEntry(key, data)

		if size+int64(len(entry)) > s.segmentSize {
			s.currentSegment++
			continue
		}

		foundCurrentSegmentFolder = false
	}

	offset, err := s.fileAdapter.Set(key, data, s.currentSegment)
	if err != nil {
		return 0, 0, fmt.Errorf("set in segment manager: %e", err)
	}

	return offset, s.currentSegment, nil
}

func (s *SegmentManager) Get(offset, segment int64) (map[string]any, error) {

	return s.fileAdapter.Get(offset, segment)
}
