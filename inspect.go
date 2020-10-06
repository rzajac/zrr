package zrr

import (
	"time"
)

// IsImmutable returns true if error err is instance of Error and is immutable.
func IsImmutable(err error) bool {
	if err, ok := err.(*Error); ok && err != nil {
		return err.imm
	}
	return false
}

// HasKey returns true if error err is instance of Error and has the key set.
func HasKey(err error, key string) bool {
	if err, ok := err.(*Error); ok && err != nil {
		_, ok := err.meta[key]
		return ok
	}
	return false
}

// HasCode returns true if error err is instance of Error and has key KCode
// which equals to one of the codes.
func HasCode(err error, codes ...string) bool {
	if err, ok := err.(*Error); ok && err != nil {
		if val, ok := err.meta[KCode]; ok {
			for _, code := range codes {
				if code == val {
					return true
				}
			}
		}
	}
	return false
}

// GetCode returns error code if error err is instance of Error.
// If error code is not set it will return empty string.
func GetCode(err error) string {
	code, _ := GetStr(err, KCode)
	return code
}

// GetStr returns the key as a string if err is an instance of Error and key
// exists. If key does not exist or it's not a string it will return
// false as the second return value.
func GetStr(err error, key string) (string, bool) {
	if err, ok := err.(*Error); ok && err != nil {
		if val, ok := err.meta[key]; ok {
			if ret, ok := val.(string); ok {
				return ret, true
			}
		}
	}
	return "", false
}

// GetInt returns the key as an integer if err is an instance of Error and key
// exists. If key does not exist or it's not an integer it will return
// false as the second return value.
func GetInt(err error, key string) (int, bool) {
	if err, ok := err.(*Error); ok && err != nil {
		if val, ok := err.meta[key]; ok {
			if ret, ok := val.(int); ok {
				return ret, true
			}
		}
	}
	return 0, false
}

// GetFloat64 returns the key as a float64 if err is an instance of Error
// and key exists. If key does not exist or it's not a float64 it will return
// false as the second return value.
func GetFloat64(err error, key string) (float64, bool) {
	if err, ok := err.(*Error); ok && err != nil {
		if val, ok := err.meta[key]; ok {
			if ret, ok := val.(float64); ok {
				return ret, true
			}
		}
	}
	return 0, false
}

// GetTime returns the key as a time.Time if err is an instance of Error
// and key exists. If key does not exist or it's not a time.Time it will return
// false as the second return value.
func GetTime(err error, key string) (time.Time, bool) {
	if err, ok := err.(*Error); ok && err != nil {
		if val, ok := err.meta[key]; ok {
			if ret, ok := val.(time.Time); ok {
				return ret, true
			}
		}
	}
	return time.Time{}, false
}

// GetBool returns the key as a boolean if err is an instance of Error and key
// exists. If key does not exist or it is not a boolean it will return
// false as the second return value.
func GetBool(err error, key string) (bool, bool) {
	if err, ok := err.(*Error); ok && err != nil {
		if val, ok := err.meta[key]; ok {
			if ret, ok := val.(bool); ok {
				return ret, true
			}
		}
	}
	return false, false
}
