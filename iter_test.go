package zrr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rzajac/zrr"
)

func Test_iter(t *testing.T) {
	// --- Given ---
	err := zrr.New("msg").Str("k0", "v0").Str("k2", "v2").Str("k1", "v1")

	// --- When ---
	keys := make([]string, 0)
	vals := make([]interface{}, 0)

	iter := err.Fields()
	for iter.Next() {
		k, v := iter.Get()
		keys = append(keys, k)
		vals = append(vals, v)
	}

	// --- Then ---
	assert.Exactly(t, []string{"k0", "k1", "k2"}, keys)
	assert.Exactly(t, []interface{}{"v0", "v1", "v2"}, vals)
}

func Test_iter_NoMetadata(t *testing.T) {
	// --- Given ---
	err := zrr.New("msg")

	// --- When ---
	keys := make([]string, 0)
	vals := make([]interface{}, 0)

	iter := err.Fields()
	for iter.Next() {
		k, v := iter.Get()
		keys = append(keys, k)
		vals = append(vals, v)
	}

	// --- Then ---
	assert.Len(t, keys, 0)
	assert.Len(t, vals, 0)
}
