package zrr

import (
	"sort"
)

// iter implements Error metadata iterator.
type iter struct {
	err  *Error   // Error to iterate over key value pairs.
	keys []string // Sorted metadata keys.
	idx  int      // Current keys slice index.
}

// newIter returns new instance of iter.
func newIter(err *Error) *iter {
	// Sort metadata keys.
	keys := make([]string, 0, len(err.meta))
	for fn := range err.meta {
		keys = append(keys, fn)
	}
	sort.Strings(keys)

	itr := &iter{
		err:  err,
		keys: keys,
		idx:  -1,
	}
	return itr
}

// Next returns true if there are more metadata keys.
// Example:
//
//    iter := err.Fields()
//    for iter.Next() {
//        k, v := iter.Get()
//        fmt.Println(k, v)
//    }
//
func (i *iter) Next() bool {
	i.idx++
	return i.idx < len(i.keys)
}

// Get returns current value of key and value from the map.
// Calling Get after Next returned false will panic.
func (i *iter) Get() (string, interface{}) {
	key := i.keys[i.idx]
	return key, i.err.meta[key]
}
