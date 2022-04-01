package localcache

import (
	"crypto/rand"
	"math"
	"math/big"
	insecurerand "math/rand"
	"os"
	"unsafe"
)

type HashFunc interface {
	Sum64(string) uint64
}

func NewDefaultHashFunc() HashFunc {
	return fnv64a{}
}
func NewHashWithDjb() HashFunc {
	return newTimes33()
}

/*****************************************************************************************************/
type fnv64a struct{}

const (
	// offset64 FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	// prime64 FNVa prime value. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	prime64 = 1099511628211
)

// Sum64 gets the string and returns its uint64 hash value.
func (f fnv64a) Sum64(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}

/*****************************************************************************************************/
func newTimes33() HashFunc {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint64
	if err != nil {
		os.Stderr.Write([]byte("WARNING: NewTimes33() failed to read from the system CSPRNG (/dev/urandom or equivalent.) Your system's security may be compromised. Continuing with an insecure seed.\n"))
		seed = uint64(insecurerand.Uint32())
	} else {
		seed = rnd.Uint64()
	}
	return djb33{seed}
}

type djb33 struct {
	seed uint64
}

func (h djb33) Sum64(k string) uint64 {
	var (
		l = uint64(len(k))
		d = 5381 + h.seed + l
		i = uint64(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint64(k[i])
			d = (d * 33) ^ uint64(k[i+1])
			d = (d * 33) ^ uint64(k[i+2])
			d = (d * 33) ^ uint64(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint64(k[i])
	case 3:
		d = (d * 33) ^ uint64(k[i])
		d = (d * 33) ^ uint64(k[i+1])
	case 4:
		d = (d * 33) ^ uint64(k[i])
		d = (d * 33) ^ uint64(k[i+1])
		d = (d * 33) ^ uint64(k[i+2])
	}
	return d ^ (d >> 16)
}

func isPowerOfTwo(number uint64) bool {
	return (number != 0) && (number&(number-1)) == 0
}
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
