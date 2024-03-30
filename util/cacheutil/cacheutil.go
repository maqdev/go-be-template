package cacheutil

import (
	"fmt"
	"strings"
	"time"
)

const KeySeparator = "|"
const RedisWriteTimeout = time.Second * 10

func Create(components ...string) string {
	return strings.Join(components, KeySeparator)
}

func Parse(key string, components ...*string) bool {
	s := strings.Split(key, KeySeparator)
	for idx, s := range s {
		if idx >= len(components) {
			return false
		}
		if components[idx] != nil {
			*components[idx] = s
		}
	}
	return len(s) == len(components)
}

// HashTag encodes key component with curly braces to specify Redis Cluster hashing slot
// https://redis.com/blog/redis-clustering-best-practices-with-keys/
func HashTag(s string) string {
	return fmt.Sprintf("{%s}", s)
}

func DecodeHashTag(s string) string {
	trimmed := strings.TrimPrefix(s, "{")
	trimmed = strings.TrimSuffix(trimmed, "}")
	if len(trimmed) == len(s)-2 {
		return trimmed
	}
	return s
}
