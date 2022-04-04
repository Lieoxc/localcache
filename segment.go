package localcache

import (
	"encoding/binary"
	"errors"
	"time"
)

const (
	chunkSize            = 16 * 1024 //每个分片16kB
	defaultIndex         = 0
	timestampSizeInBytes = 8
	hashSizeInBytes      = 8                                                       // Number of bytes used for hash
	keySizeInBytes       = 2                                                       // Number of bytes used for size of entry key
	headersSizeInBytes   = timestampSizeInBytes + hashSizeInBytes + keySizeInBytes // Number of bytes used for all headers

)

type segment struct {
	hashmap map[uint64]uint32
	chunks  [][]byte
	index   uint64
}

func newSegment(bytes uint64) *segment {
	capacity := (bytes + chunkSize - 1) / chunkSize
	chunks := make([][]byte, capacity)

	return &segment{
		chunks:  chunks,
		hashmap: make(map[uint64]uint32),
		index:   defaultIndex,
	}
}
func (s *segment) clean(timestamp int64) {
	count := []int{}
	for i := range count {
		entry := s.chunks[i]
		if entry == nil {
			continue
		}
		expireTime := readExpireAtFromEntry(entry)
		if timestamp-int64(expireTime) >= 0 {
			hash := readHashFromEntry(entry)
			s.removeChunks(i)
			delete(s.hashmap, hash)
		}
	}
}
func (s *segment) set(key string, hashKey uint64, value []byte, expireTime time.Duration) error {
	if index, ok := s.hashmap[hashKey]; ok {
		s.removeChunks(int(index))
		delete(s.hashmap, hashKey)
	}
	expireAt := time.Now().Add(expireTime).Unix()
	entry := wrapEntry(expireAt, key, hashKey, value)
	index, err := s.push(entry)
	if err == nil {
		s.hashmap[hashKey] = uint32(index)
	}
	return nil
}

func (s *segment) get(key string, hashKey uint64) ([]byte, error) {
	entry, err := s.getEntry(key, hashKey)
	if err != nil {
		return nil, err
	}
	res := readEntry(entry)
	if res == nil {
		return nil, ErrEntryNotFound
	}
	return res, nil
}
func (s *segment) del(hashKey uint64) error {
	index, ok := s.hashmap[hashKey]
	if !ok {
		return nil
	}
	s.removeChunks(int(index))
	delete(s.hashmap, hashKey)
	return nil
}
func (s *segment) len() int {
	return len(s.hashmap)
}
func (s *segment) removeChunks(index int) {
	s.chunks[index] = nil
}
func (s *segment) push(data []byte) (uint64, error) {
	dataLen := len(data)
	dataIndex := s.index
	s.chunks[s.index] = make([]byte, dataLen)
	copy(s.chunks[s.index][0:], data[:dataLen])
	s.index++
	return dataIndex, nil
}
func (s *segment) getEntry(key string, hashKey uint64) ([]byte, error) {

	index, ok := s.hashmap[hashKey]
	if !ok {
		return nil, ErrEntryNotFound
	}
	entry := s.chunks[index]
	if entry == nil {
		return nil, errors.New("entry is nil")
	}
	if entryKey := readKeyFromEntry(entry); entryKey == key {
		return entry, nil
	}
	return nil, ErrkeyERR
}

func readKeyFromEntry(data []byte) string {
	length := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, length)
	//拷贝key
	copy(dst, data[headersSizeInBytes:headersSizeInBytes+length])
	return bytesToString(dst)
}
func readHashFromEntry(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data[timestampSizeInBytes:])
}

func readEntry(data []byte) []byte {
	length := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, len(data)-int(length+headersSizeInBytes))
	//拷贝数据段
	copy(dst, data[headersSizeInBytes+length:])

	return dst
}

func readExpireAtFromEntry(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func wrapEntry(timestamp int64, key string, hashKey uint64, value []byte) []byte {
	keyLen := len(key)
	blobLength := len(value) + keyLen + headersSizeInBytes
	blob := make([]byte, blobLength)
	binary.LittleEndian.PutUint64(blob, uint64(timestamp))
	binary.LittleEndian.PutUint64(blob[timestampSizeInBytes:], hashKey)
	binary.LittleEndian.PutUint16(blob[timestampSizeInBytes+hashSizeInBytes:], uint16(keyLen))
	copy(blob[headersSizeInBytes:], key)
	copy(blob[headersSizeInBytes+keyLen:], value)

	return blob[:blobLength]
}
