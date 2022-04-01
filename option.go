package localcache

import "time"

const (
	defaultBucketCount = 256
	defaultMaxBytes    = 512 * 1024 * 1024 // 512M
	defaultCleanTIme   = time.Minute * 10
)

type options struct {
	hashFunc       HashFunc
	bucketCount    uint64
	maxBytes       uint64
	cleanTime      time.Duration
	cleanupEnabled bool
}
type Opt func(options *options)

func SetHashFunc(hashFunc HashFunc) Opt {
	return func(opt *options) {
		opt.hashFunc = hashFunc
	}
}

func SetShardCount(count uint64) Opt {
	return func(opt *options) {
		opt.bucketCount = count
	}
}

func SetMaxBytes(maxBytes uint64) Opt {
	return func(opt *options) {
		opt.maxBytes = maxBytes
	}
}

func SetCleanTime(time time.Duration) Opt {
	return func(opt *options) {
		opt.cleanTime = time
	}
}

func SetCleanupEnabled(enabled bool) Opt {
	return func(opt *options) {
		opt.cleanupEnabled = enabled
	}
}
